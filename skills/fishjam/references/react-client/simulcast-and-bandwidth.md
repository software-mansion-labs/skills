# Simulcast & Bandwidth

Simulcast lets a publisher send **multiple quality layers** of the same video track. Subscribers choose which layer to receive. The benefit: a grid showing 12 tiles can request `LOW` from every peer to save bandwidth, while the spotlit speaker gets `HIGH`.

## Enabling simulcast on publish

Pass `videoConfig` to `FishjamProvider`:

```tsx
import { FishjamProvider, Variant } from '@fishjam-cloud/react-client';

<FishjamProvider
  fishjamId={...}
  videoConfig={{
    sentQualities: [Variant.VARIANT_LOW, Variant.VARIANT_MEDIUM, Variant.VARIANT_HIGH],
  }}
>
```

`StreamConfig` shape:

```ts
type StreamConfig = { sentQualities?: Variant[] | false };
```

- `sentQualities: Variant[]` — which simulcast layers to encode and publish. The default when starting a track is `[VARIANT_LOW, VARIANT_MEDIUM, VARIANT_HIGH]` (i.e. full simulcast).
- `sentQualities: false` — disable simulcast and publish a single (non-layered) encoding.
- `sentQualities: undefined` — fall back to the default of all three layers.

Drop layers if you want simpler simulcast — e.g. `[Variant.VARIANT_LOW, Variant.VARIANT_HIGH]` skips medium.

Per-layer bandwidth caps live on `bandwidthLimits` (below), not on `StreamConfig`.

## Variant enum

```ts
import { Variant } from '@fishjam-cloud/react-client';
// Variant.VARIANT_LOW, Variant.VARIANT_MEDIUM, Variant.VARIANT_HIGH
```

Re-exported from `@fishjam-cloud/ts-client` (originally from the protobuf shared types). Use these constants rather than string literals or unprefixed names.

## Asking for a specific layer per track

```tsx
import { usePeers, Variant } from '@fishjam-cloud/react-client';

const { remotePeers, setReceivedTracksQuality } = usePeers();

// One track:
peer.cameraTrack?.setReceivedQuality(Variant.VARIANT_LOW);

// Bulk:
const allRemoteTrackIds = remotePeers.flatMap((p) => p.tracks.map((t) => t.trackId));
setReceivedTracksQuality(allRemoteTrackIds, Variant.VARIANT_LOW);
```

Behavior:

- The request is best-effort — Fishjam delivers the closest available layer.
- If the publisher didn't enable simulcast, layer selection is a no-op (you always get the only encoding).
- The SFU does not auto-select a layer based on element size — call `setReceivedQuality` explicitly when you want anything other than the default (typically the highest available).

## Bandwidth limits

`bandwidthLimits` on the provider caps bandwidth per video stream:

```tsx
import { FishjamProvider, Variant } from '@fishjam-cloud/react-client';

<FishjamProvider
  fishjamId={...}
  bandwidthLimits={{
    singleStream: 1500,    // kbps cap when not using simulcast (one encoding)
    simulcast: {           // kbps cap per simulcast layer
      [Variant.VARIANT_LOW]: 150,
      [Variant.VARIANT_MEDIUM]: 500,
      [Variant.VARIANT_HIGH]: 1500,
    },
  }}
>
```

Shape (from `packages/react-client/src/types/public.ts`):

```ts
type BandwidthLimits = {
  singleStream: number;                        // kbps
  simulcast: {
    [Variant.VARIANT_LOW]: number;
    [Variant.VARIANT_MEDIUM]: number;
    [Variant.VARIANT_HIGH]: number;
  };
};
```

A value of `0` means "no cap" — the SDK default. Lower these only when targeting low-bandwidth networks. There is no audio entry; audio bandwidth is controlled by the codec, not by this object.

## Decision guide

- **All peers see all peers in a grid** → enable simulcast, request `VARIANT_LOW` everywhere except focused tiles.
- **One presenter + 50 listeners (livestream)** → use `livestream` room type, no simulcast needed (Fishjam handles fan-out).
- **High-quality 1:1** → simulcast not needed; set `videoConfig={{ sentQualities: false }}` and bump `bandwidthLimits.singleStream` to 2500 if you want HD.
- **Bandwidth-constrained mobile** → cap `bandwidthLimits.singleStream` to 800, or cap the per-layer values under `bandwidthLimits.simulcast` and request `Variant.VARIANT_LOW` on the receive side.

## Audio simulcast?

Rare. The provider accepts `audioConfig` with the same shape, but audio simulcast is unusual — the codec doesn't benefit the way video does. Leave it `undefined` unless you have a specific need.

## Sources

- `@fishjam-cloud/react-client` → `packages/react-client/src/FishjamProvider.tsx`
- `packages/react-client/src/types/public.ts` (`StreamConfig`, `BandwidthLimits`, `SimulcastBandwidthLimits`)
- `packages/react-client/src/utils/bandwidth.ts` (`mergeWithDefaultBandwitdthLimits`, `ALL_VARIANTS_SIMULCAST`)
- `packages/react-client/src/hooks/internal/useTrackManager.ts` (default `sentQualities`)
- `packages/protobufs/fishjam/media_events/shared.ts` (`Variant` enum)
- `packages/react-client/src/hooks/usePeers.ts` (`setReceivedTracksQuality`)
- <https://fishjam.swmansion.com/docs/how-to/client/simulcast>
- <https://fishjam.swmansion.com/docs/explanation/simulcast>
