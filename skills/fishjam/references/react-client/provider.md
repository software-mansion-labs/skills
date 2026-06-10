# FishjamProvider

The single React context provider. Wrap your app in it once.

Source: `@fishjam-cloud/react-client` → `packages/react-client/src/FishjamProvider.tsx`.

## Minimal

```tsx
import { FishjamProvider } from '@fishjam-cloud/react-client';

<FishjamProvider fishjamId={import.meta.env.VITE_FISHJAM_ID}>
  <YourApp />
</FishjamProvider>
```

## All props

| Prop | Type | Default | Notes |
|---|---|---|---|
| `fishjamId` | `string` | **required** | Your Fishjam instance ID from the Dashboard. Read by `useConnection` to construct the WebSocket URL. |
| `reconnect` | `ReconnectConfig \| boolean` | `true` | Auto-reconnect policy. `true` = default policy. `false` = no auto-reconnect. Object form: `{ maxAttempts, initialDelay, delay, addTracksOnReconnect }` — see `@fishjam-cloud/ts-client` types. |
| `constraints` | `{ audio?, video? }` | `{ video: VIDEO_TRACK_CONSTRAINTS, audio: true }` | `MediaStreamConstraints.video` / `.audio` defaults used by `useCamera` / `useMicrophone` before per-device overrides kick in. |
| `persistLastDevice` | `boolean \| PersistLastDeviceHandlers` | `undefined` (= persist via localStorage when provider mounts in browser) | Remember the last-used camera/mic across page loads. `false` = off. Provide handlers `{ getLastDevice, saveLastDevice }` for custom storage (e.g. cookies, IndexedDB). |
| `bandwidthLimits` | `Partial<BandwidthLimits>` | `0` per layer = unlimited | Per-encoding bandwidth caps; merged with `{ singleStream: 0, simulcast: { LOW: 0, MEDIUM: 0, HIGH: 0 } }` via `mergeWithDefaultBandwitdthLimits`. Default is `0` everywhere — no cap. |
| `videoConfig` | `StreamConfig` | `undefined` | Enables video simulcast + which `Variant` layers to publish. See `simulcast-and-bandwidth.md`. |
| `audioConfig` | `StreamConfig` | `undefined` | Same for audio (rarely needed; audio simulcast is uncommon). |
| `debug` | `boolean` | `false` | Toggles ts-client debug logs in console. |
| `fishjamClient` | `FishjamClient` (from `@fishjam-cloud/ts-client`) | (constructed internally) | Override the underlying client. Useful for tests or for sharing one client with non-React code. **When you pass `fishjamClient`, the provider ignores `reconnect` and `debug` — those are constructor args for the internally-built client only.** |

## ReconnectConfig (boolean shorthand)

```tsx
<FishjamProvider fishjamId={...} reconnect={false}>          {/* never reconnect */}
<FishjamProvider fishjamId={...} reconnect={true}>           {/* default policy */}
<FishjamProvider fishjamId={...} reconnect={{                /* explicit */
  maxAttempts: 5,            // default: 3
  initialDelay: 1000,        // ms before first retry; default: 500
  delay: 2000,               // ms between subsequent retries; default: 500
  addTracksOnReconnect: true, // re-add published tracks after reconnect; default: false
}}>
```

Status of any in-flight reconnect is exposed via `useConnection().reconnectionStatus`.

## Constraints — example: HD video

```tsx
<FishjamProvider
  fishjamId={...}
  constraints={{
    video: { width: { ideal: 1280 }, height: { ideal: 720 }, frameRate: { ideal: 30 } },
    audio: { echoCancellation: true, noiseSuppression: true },
  }}
>
```

These defaults apply when the user hasn't picked a specific device. `useCamera().selectCamera(deviceId)` overrides.

## persistLastDevice — custom storage

```tsx
const handlers = {
  getLastDevice: async () => { /* return DeviceLabels or null */ },
  saveLastDevice: async (device) => { /* persist */ },
};

<FishjamProvider fishjamId={...} persistLastDevice={handlers}>
```

Default implementation uses `localStorage`; replace if you need cross-domain or SSR-safe behavior.

## bandwidthLimits

```tsx
import { FishjamProvider, Variant } from '@fishjam-cloud/react-client';

<FishjamProvider
  fishjamId={...}
  bandwidthLimits={{
    singleStream: 1500,    // kbps cap for non-simulcast video
    simulcast: {
      [Variant.VARIANT_LOW]: 150,
      [Variant.VARIANT_MEDIUM]: 500,
      [Variant.VARIANT_HIGH]: 1500,
    },
  }}
>
```

If you're shipping video to bandwidth-constrained clients, lower these to avoid stalls.

## fishjamClient override

Construct a `FishjamClient` (from `@fishjam-cloud/ts-client`) yourself and pass it in. Useful when:

- Sharing one client between React and a worker.
- Driving the same client from a non-React harness in tests.
- Using a custom subclass for instrumentation.

```tsx
import { FishjamClient } from '@fishjam-cloud/ts-client';

const client = new FishjamClient({ reconnect: true, debug: false });

<FishjamProvider fishjamId={...} fishjamClient={client}>
```

If omitted, `FishjamProvider` constructs one internally from `reconnect` and `debug` props.

## Sources

- `@fishjam-cloud/react-client` → `packages/react-client/src/FishjamProvider.tsx`
- `packages/react-client/src/types/public.ts` (`BandwidthLimits`, `SimulcastBandwidthLimits`, `StreamConfig`)
- `packages/ts-client/src/reconnection.ts` (`ReconnectConfig`)
- <https://fishjam.swmansion.com/docs/how-to/client/installation>
- <https://fishjam.swmansion.com/docs/how-to/client/reconnection-handling>
