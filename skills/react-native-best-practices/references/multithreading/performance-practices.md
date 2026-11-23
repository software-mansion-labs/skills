# Worklets Performance Practices

Use these practices when writing or reviewing code that crosses React Native, UI, or Worker Runtime boundaries. They minimize runtime creation, worklet serialization, data transfer, scheduling, and thread-blocking overhead.

Source: [React Native Worklets: Good practices](https://docs.swmansion.com/react-native-worklets/docs/guides/good-practices/)

## Reuse Worker Runtimes

Create a Worker Runtime once at app startup or in module scope and retain its reference. Each `createWorkletRuntime` call creates a JavaScript runtime and, with the default queue, a thread.

Do not create a runtime in a component body or frequently called function. Create a fresh runtime only when the task specifically needs an isolated, clean environment.

```tsx
const workerRuntime = createWorkletRuntime({ name: 'worker' });

function processData(data: Data) {
  'worklet';
  return transform(data);
}

function processInBackground(data: Data) {
  scheduleOnRuntime(workerRuntime, processData, data);
}
```

## Reuse Worklet Function Objects

Define repeatedly scheduled worklet functions once in module scope. If the function must live in a component, stabilize it with `useCallback`, but prefer passing changing values as arguments so dependencies do not recreate it.

Worklets caches serialization by function-object identity and retains the resulting function on the target runtime. An inline arrow or component-local function created on each call or render cannot reuse those cached forms, especially when its closure is large.

```tsx
function logValue(value: number) {
  'worklet';
  console.log(value);
}

function onPress(value: number) {
  scheduleOnUI(logValue, value);
}
```

## Pass Changing Values as Arguments

Prefer direct arguments to captured values. This keeps the worklet function reusable and transfers the value at the time of the call.

```tsx
function logTitle(title: string) {
  'worklet';
  console.log(title);
}

function Header({ title }: { title: string }) {
  const onPress = () => scheduleOnUI(logTitle, title);
  return <Button onPress={onPress} title={title} />;
}
```

Do not introduce a wrapper worklet solely to capture arguments and forward them to another worklet. Pass the target function and its arguments directly:

```tsx
function schedule(arg1: string, arg2: number) {
  scheduleOnUI(foo, arg1, arg2);
}
```

Keep a wrapper only when it performs additional work on the target runtime.

## Transfer Only Required Data

Pass or capture only the fields the worklet needs. Accessing one property from a captured object serializes the whole object graph, including nested values. In development builds, serialization also makes the original object immutable.

```tsx
function reportStatus(status: string) {
  'worklet';
  console.log(status);
}

scheduleOnUI(reportStatus, report.status);
```

Read the property at the scheduling call when the worklet needs its current value. Capture a whole object only when the target operation actually needs the whole object.

## Use the Target Runtime's Event Loop

For delayed or repeating work, install the timer on the runtime that performs the work. Do not keep a timer on the RN Runtime that sends a new cross-runtime job on every tick.

Worker Runtimes created with the default `enableEventLoop` setting provide timers, `setImmediate`, `requestAnimationFrame`, and `queueMicrotask`. The UI Runtime also provides timers.

```tsx
function tick(): boolean {
  'worklet';
  return updateWorkerState();
}

scheduleOnRuntime(workerRuntime, () => {
  'worklet';
  const intervalId = setInterval(() => {
    if (tick()) {
      clearInterval(intervalId);
    }
  }, 100);
});
```

Keep a repeating timer's identifier and call `clearInterval` on the same runtime when the work completes or its owner is disposed. Use an RN Runtime timer only when the timing decision itself belongs to RN-side state or APIs.

## Prefer Async Results

Use `runOnUIAsync` or `runOnRuntimeAsync` when returning a value unless the caller genuinely must have the result synchronously. Synchronous variants block the calling thread until the target finishes and may preempt work on the target runtime.

```tsx
async function onPress() {
  const layout = await runOnUIAsync(measureLayout);
  displayLayout(layout);
}
```

Treat `runOnUISync` and `runOnRuntimeSync` as deliberate latency trade-offs. Keep them only when synchronous ordering is required and the work is bounded enough to justify blocking the caller.

## Review Checklist

When performance or multithreading overhead is in scope, check whether:

1. Worker Runtimes are long-lived unless isolation requires a fresh environment.
2. Repeatedly scheduled worklets keep the same function identity.
3. Changing values are direct arguments rather than captured closures.
4. Wrapper worklets perform real target-runtime work rather than only forwarding arguments.
5. Only required fields cross the runtime boundary.
6. Delayed and repeating work uses the target runtime's event loop and is cleaned up there.
7. Calls that return values are asynchronous unless synchronous ordering is necessary.
