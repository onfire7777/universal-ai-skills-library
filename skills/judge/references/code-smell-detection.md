# Code Smell & Design Anti-Pattern Detection

> レビュー時に検出すべき構造的コードスメル、設計アンチパターン、検出ヒューリスティクス

## 1. コードスメル検出の役割

```
Judge の検出レイヤー:
  Layer 1: バグパターン    → bug-patterns.md（既存）
  Layer 2: 一貫性問題      → consistency-patterns.md（既存）
  Layer 3: コードスメル    → 本ファイル（新規）
  Layer 4: テスト品質      → test-quality-patterns.md（既存）

コードスメル = 即座のバグではないが、
  メンテナンス性・拡張性・品質を劣化させる構造的問題
```

---

## 2. 検出すべきコードスメル Top 10

### SMELL-001: God Class / God Object

```
定義: 過剰な責務を持つ単一クラス/モジュール
検出: メソッド数 > 20 or 行数 > 500 or 依存クラス数 > 10

影響: 変更の波及範囲が予測不能、テスト困難
severity: MEDIUM
route: → Zen（リファクタリング）
```

```typescript
// ❌ God Class
class UserManager {
  createUser() { ... }
  deleteUser() { ... }
  sendEmail() { ... }      // メール送信?
  generateReport() { ... }  // レポート生成?
  processPayment() { ... }  // 決済処理?
  validateAddress() { ... } // 住所検証?
}

// ✅ 単一責任に分割
class UserService { createUser(); deleteUser(); }
class EmailService { sendEmail(); }
class PaymentService { processPayment(); }
```

### SMELL-002: Spaghetti Code

```
定義: 制御フローが複雑に絡み合った構造
検出: 循環的複雑度 > 15 or ネスト深度 > 4 or goto/break の多用

影響: デバッグ困難、バグ混入リスク高
severity: MEDIUM
route: → Zen（リファクタリング）
```

### SMELL-003: Primitive Obsession

```
定義: 複雑な概念を基本型（string, number）で表現
検出: 同じ型の引数が3つ以上 or ドメイン概念が string で表現

影響: 型安全性低下、バリデーション漏れ
severity: LOW
route: → Zen（Value Object 導入）
```

```typescript
// ❌ Primitive Obsession
function createOrder(
  userId: string,
  productId: string,
  quantity: number,
  price: number,
  currency: string,
) { ... }

// ✅ ドメイン型を使用
function createOrder(
  userId: UserId,
  productId: ProductId,
  quantity: Quantity,
  price: Money,
) { ... }
```

### SMELL-004: Shotgun Surgery

```
定義: 1つの変更が多数のファイルに波及
検出: 1つの概念変更で 5+ ファイル変更 or 同一修正の散在

影響: 修正漏れリスク、変更コスト増大
severity: MEDIUM
route: → Atlas（アーキテクチャ改善）
```

### SMELL-005: Feature Envy

```
定義: あるクラスが別クラスのデータに過度にアクセス
検出: 他オブジェクトのプロパティへのチェーンアクセス（a.b.c.d）

影響: カプセル化違反、結合度増大
severity: LOW
route: → Zen（メソッド移動）
```

### SMELL-006: Long Method / Function

```
定義: 1つの関数が過度に長い
検出: 行数 > 50 or パラメータ数 > 5 or 処理ステップ > 3

影響: 理解困難、テスト困難、再利用不可
severity: LOW
route: → Zen（関数抽出）
```

### SMELL-007: Dead Code

```
定義: 実行されないコード、到達不能なコード
検出: 未使用の関数/変数/import、コメントアウトされたコード

影響: コードベース肥大化、誤解の原因
severity: INFO
route: → Sweep（不要コード削除）
```

### SMELL-008: Magic Numbers / Strings

```
定義: 意味が不明なリテラル値の直接使用
検出: コード中の数値/文字列リテラル（0, 1, -1, "" を除く）

影響: 意図不明、変更時の修正漏れ
severity: INFO
route: → Zen（定数化）
```

### SMELL-009: Duplicated Logic

```
定義: 同一/類似ロジックの複数箇所での繰り返し
検出: 3行以上の類似コードブロック x 2+ 箇所

影響: 修正時の一貫性欠如、バグの温床
severity: LOW
route: → Zen（共通化）
```

### SMELL-010: Inappropriate Intimacy

```
定義: クラス間の過度な密結合
検出: 友達クラスの private メンバーへのアクセス、循環参照

影響: 変更の波及、独立テスト困難
severity: MEDIUM
route: → Atlas（依存関係整理）
```

---

## 3. 検出しきい値マトリクス

| スメル | 検出指標 | しきい値 | Severity |
|--------|---------|---------|----------|
| God Class | メソッド数 | > 20 | MEDIUM |
| God Class | クラス行数 | > 500 | MEDIUM |
| Spaghetti | 循環的複雑度 | > 15 | MEDIUM |
| Spaghetti | ネスト深度 | > 4 | MEDIUM |
| Long Function | 関数行数 | > 50 | LOW |
| Long Function | パラメータ数 | > 5 | LOW |
| Duplication | 類似ブロック | 3行+ x 2箇所+ | LOW |
| Dead Code | 未使用 export | 任意 | INFO |
| Magic Number | リテラル値 | context 依存 | INFO |

---

## 4. Judge レポートでのスメル報告

### 報告フォーマット

```
## Code Smell Findings

| ID | Type | File:Line | Description | Severity | Route |
|----|------|-----------|-------------|----------|-------|
| SMELL-001 | God Class | src/services/UserManager.ts | 35 methods, 890 LOC | MEDIUM | → Zen |
| SMELL-006 | Long Function | src/utils/transform.ts:45 | 120 LOC, 8 params | LOW | → Zen |
| SMELL-007 | Dead Code | src/legacy/old-helper.ts | Unused export (3 functions) | INFO | → Sweep |
```

### 報告ポリシー

```
Judge のスメル報告ルール:
  1. バグパターンを優先（スメルは補助的情報）
  2. MEDIUM 以上のスメルのみ主要レポートに含める
  3. INFO/LOW は "Additional Observations" セクションに
  4. 新規導入のスメルのみ報告（既存コードは対象外）
  5. route 先エージェントを必ず指定
```

---

## 5. フレームワーク固有のスメル

### React

```
- Prop Drilling（3層以上のプロップ転送）→ Context/State 管理
- useEffect 依存配列の不備 → bug-patterns.md で検出
- コンポーネント肥大化（300行+）→ 分割推奨
```

### Express / API

```
- Fat Controller（ルートハンドラにビジネスロジック）→ Service 層分離
- エラーハンドリングの不統一 → consistency-patterns.md で検出
- ミドルウェアチェーンの過度な複雑さ → 簡素化推奨
```

### TypeScript

```
- any 型の多用 → 型定義化（→ Quill）
- type assertion（as）の多用 → 型設計の見直し
- enum vs union type の不統一 → consistency-patterns.md で検出
```

**Source:** [CodeRabbit: 5 Code Review Anti-Patterns](https://www.coderabbit.ai/blog/5-code-review-anti-patterns-you-can-eliminate-with-ai) · [DZone: Code Review Patterns and Anti-Patterns](https://dzone.com/refcardz/code-review-patterns-and-anti-patterns) · [HackerNoon: Code Review Anti-Patterns](https://hackernoon.com/code-review-anti-patterns-how-to-stop-nitpicking-syntax-and-start-improving-architecture)
