# Screen capture — Android

MediaProjection with three moving parts: the consent intent, a `mediaProjection`-typed foreground service, and `ScreenCapture` itself. Video-only — unlike iOS there is no app-audio source; add a `MicrophoneCapture` track if the broadcast needs sound.

Manifest:

```xml
<uses-permission android:name="android.permission.FOREGROUND_SERVICE" />
<uses-permission android:name="android.permission.FOREGROUND_SERVICE_MEDIA_PROJECTION" />
<service android:name=".ScreenCaptureService"
         android:foregroundServiceType="mediaProjection"
         android:exported="false" />
```

Flow (Android 14+ requires the foreground service to start **from inside the consent result callback**):

```kotlin
val launcher = rememberLauncherForActivityResult(StartActivityForResult()) { result ->
    if (result.resultCode == Activity.RESULT_OK) {
        // MUST happen here on Android 14+
        context.startForegroundService(Intent(context, ScreenCaptureService::class.java))
        viewModel.publishScreen(result.data!!, result.resultCode)
    }
}
launcher.launch(context.getSystemService(MediaProjectionManager::class.java).createScreenCaptureIntent())
```

```kotlin
// In the service: startForeground(id, notification, ServiceInfo.FOREGROUND_SERVICE_TYPE_MEDIA_PROJECTION)
// In the view model, once the service reports it is running:
val screen = ScreenCapture(intent = data, resultCode = resultCode, width = 1280, height = 720, frameRate = 30)
screen.start(context) // suspend; foreground service must already be running
publisher.addVideoTrack(name = "screen", source = screen, config = videoConfig)
```

Wait for the service to actually be in the foreground before `screen.start()` — the demo's service exposes `isRunning: StateFlow<Boolean>` plus a `suspend fun awaitStarted()` and the publish flow wraps it in `withTimeout(5_000)`. Stop order: `publisher.stop()` → `screen.stop()` → stop the service.
