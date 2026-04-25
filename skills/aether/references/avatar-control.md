# Avatar Control

Purpose: Read this when choosing `Live2D` vs `VRM`, defining avatar-control contracts, mapping emotions to parameters, or designing idle motion.

## Contents

- [Framework comparison](#framework-comparison)
- [Live2D Cubism SDK for Web](#live2d-cubism-sdk-for-web)
- [VRM](#vrm)
- [Emotion to expression mapping](#emotion-to-expression-mapping)
- [Idle motion design](#idle-motion-design)

## Framework Comparison

| Feature | Live2D Cubism SDK for Web | VRM (`@pixiv/three-vrm`) |
|---------|---------------------------|--------------------------|
| Model format | `.moc3` + `.model3.json` | `.vrm` (glTF extension) |
| Rendering | WebGL (Cubism Core) | Three.js + WebGL |
| Avatar style | 2D layered mesh deformation | 3D skeletal animation |
| Expression system | Parameter-based (`0.0-1.0`) | BlendShape preset + custom |
| Lip sync | Mouth-open parameter | Viseme BlendShapes (`aa/ih/ou/ee/oh`) |
| Physics | Built-in physics | Spring bone extension |
| Idle motion | Motion files (`.motion3.json`) | VRM Animation / custom |
| Community fit | VTuber standard in JP | Open cross-platform standard |
| License | Cubism SDK license | MIT (`three-vrm`) |
| Recommended for | Traditional 2D VTuber | Full-body or 3D VTuber |

## Live2D Cubism SDK for Web

### Parameter Control

Live2D models are controlled by named float parameters.

| Parameter ID | Range | Description |
|-------------|-------|-------------|
| `ParamAngleX` | `-30` to `30` | Head rotation left / right |
| `ParamAngleY` | `-30` to `30` | Head rotation up / down |
| `ParamAngleZ` | `-30` to `30` | Head tilt |
| `ParamBodyAngleX` | `-10` to `10` | Body rotation |
| `ParamEyeLOpen` | `0` to `1` | Left eye open |
| `ParamEyeROpen` | `0` to `1` | Right eye open |
| `ParamEyeBallX` | `-1` to `1` | Eye gaze horizontal |
| `ParamEyeBallY` | `-1` to `1` | Eye gaze vertical |
| `ParamBrowLY` | `-1` to `1` | Left eyebrow position |
| `ParamBrowRY` | `-1` to `1` | Right eyebrow position |
| `ParamMouthOpenY` | `0` to `1` | Mouth open amount |
| `ParamMouthForm` | `-1` to `1` | Mouth shape (`-1=frown`, `1=smile`) |

### Setup (Cubism SDK for Web)

```typescript
import { CubismFramework } from '@cubism-sdk/framework';

CubismFramework.startUp();
CubismFramework.initialize();

const model = await loadModel('model.model3.json');

model.setParameterValueById('ParamMouthOpenY', 0.8);
model.setParameterValueById('ParamMouthForm', 0.5);

function animate() {
  model.update();
  renderer.draw(model);
  requestAnimationFrame(animate);
}
```

### Motion System

```typescript
const motion = await loadMotion('idle.motion3.json');
model.startMotion(motion, {
  group: 'idle',
  priority: 2,  // 1=idle, 2=normal, 3=force
  onFinish: () => { /* return to idle */ },
});

// Motion priority:
// 1 (Idle)    continuous background motion
// 2 (Normal)  event-triggered motion
// 3 (Force)   overrides everything
```

## VRM

### BlendShape Control

| Preset | Meaning | Typical use |
|--------|---------|-------------|
| `happy` | Smile | Joy |
| `angry` | Anger | Angry state |
| `sad` | Sadness | Sad state |
| `relaxed` | Relaxed | Neutral / calm |
| `surprised` | Surprise | Surprise |
| `aa` | Mouth shape A | Lip sync vowel A |
| `ih` | Mouth shape I | Lip sync vowel I |
| `ou` | Mouth shape U | Lip sync vowel U |
| `ee` | Mouth shape E | Lip sync vowel E |
| `oh` | Mouth shape O | Lip sync vowel O |
| `blink` | Blink | Auto-blink |
| `blinkLeft` | Left wink | Expression |
| `blinkRight` | Right wink | Expression |
| `lookUp/Down/Left/Right` | Gaze | Eye direction |

### Setup (`three-vrm`)

```typescript
import * as THREE from 'three';
import { GLTFLoader } from 'three/addons/loaders/GLTFLoader.js';
import { VRMLoaderPlugin, VRM } from '@pixiv/three-vrm';

const loader = new GLTFLoader();
loader.register((parser) => new VRMLoaderPlugin(parser));

const gltf = await loader.loadAsync('model.vrm');
const vrm: VRM = gltf.userData.vrm;
scene.add(vrm.scene);

vrm.expressionManager?.setValue('happy', 0.8);
vrm.expressionManager?.setValue('aa', 1.0);

function animate() {
  vrm.update(clock.getDelta());
  renderer.render(scene, camera);
  requestAnimationFrame(animate);
}
```

## Emotion to Expression Mapping

### Live2D Expression Map

| Emotion | ParamEyeLOpen | ParamEyeROpen | ParamBrowLY | ParamBrowRY | ParamMouthForm | ParamMouthOpenY |
|---------|---------------|---------------|-------------|-------------|----------------|-----------------|
| Neutral | `1.0` | `1.0` | `0.0` | `0.0` | `0.0` | `0.0` |
| Joy | `0.8` | `0.8` | `0.3` | `0.3` | `0.8` | `0.3` |
| Sad | `0.6` | `0.6` | `-0.5` | `-0.5` | `-0.5` | `0.0` |
| Angry | `0.9` | `0.9` | `-0.7` | `-0.7` | `-0.3` | `0.2` |
| Surprised | `1.0` | `1.0` | `0.8` | `0.8` | `0.0` | `0.7` |
| Thinking | `0.7` | `0.9` | `0.2` | `-0.2` | `0.0` | `0.0` |

### VRM Expression Map

| Emotion | happy | angry | sad | surprised | relaxed | Custom blendshapes |
|---------|-------|-------|-----|-----------|---------|--------------------|
| Neutral | `0.0` | `0.0` | `0.0` | `0.0` | `0.3` | — |
| Joy | `0.8` | `0.0` | `0.0` | `0.0` | `0.2` | — |
| Sad | `0.0` | `0.0` | `0.7` | `0.0` | `0.0` | — |
| Angry | `0.0` | `0.8` | `0.0` | `0.0` | `0.0` | — |
| Surprised | `0.0` | `0.0` | `0.0` | `0.9` | `0.0` | — |
| Thinking | `0.0` | `0.0` | `0.0` | `0.0` | `0.5` | `lookUp: 0.3` |

### Expression Transition

Do not switch expressions instantly; interpolate them.

```typescript
class ExpressionTransitioner {
  private current: Record<string, number> = {};
  private target: Record<string, number> = {};
  private transitionSpeed = 0.1;

  setTarget(expression: Record<string, number>): void {
    this.target = { ...expression };
  }

  update(deltaTime: number): Record<string, number> {
    for (const [key, targetValue] of Object.entries(this.target)) {
      const current = this.current[key] ?? 0;
      const diff = targetValue - current;
      this.current[key] = current + diff * Math.min(this.transitionSpeed * deltaTime * 60, 1);
    }
    return { ...this.current };
  }
}
```

## Idle Motion Design

Use idle motion whenever no active speech or explicit reaction is running.

### Idle Behaviors

| Behavior | Frequency | Duration | Parameters |
|----------|-----------|----------|------------|
| Auto-blink | `3-5s` interval | `150ms` | `EyeLOpen/ROpen: 1→0→1` |
| Breathing | Continuous | `4s` cycle | `BodyAngleX: ±1`, scale `±0.02` |
| Head sway | Continuous | `6-8s` cycle | `AngleX: ±3`, `AngleY: ±2` |
| Eye wander | `5-10s` interval | `1-2s` | `EyeBallX/Y: random ±0.3` |
| Micro expression | `15-30s` interval | `2-3s` | `MouthForm: 0→0.2→0` |

### Implementation Pattern

```typescript
class IdleAnimator {
  private time = 0;

  update(deltaTime: number): Record<string, number> {
    this.time += deltaTime;

    return {
      bodyScale: 1 + Math.sin(this.time * Math.PI * 0.5) * 0.02,
      headAngleX: Math.sin(this.time * Math.PI * 2 / 7) * 3,
      headAngleY: Math.sin(this.time * Math.PI * 2 / 9) * 2,
      eyeOpen: this.calculateBlink(this.time),
    };
  }

  private calculateBlink(time: number): number {
    const blinkInterval = 4;
    const blinkDuration = 0.15;
    const t = time % blinkInterval;
    if (t < blinkDuration) {
      return t < blinkDuration / 2
        ? 1 - (t / (blinkDuration / 2))
        : (t - blinkDuration / 2) / (blinkDuration / 2);
    }
    return 1;
  }
}
```
