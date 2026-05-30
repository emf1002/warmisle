# Directory Structure Cleanup Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Clean up repository by removing design artifacts, consolidating build scripts, and reorganizing frontend test utilities.

**Architecture:** This is a pure file reorganization task — no code logic changes. Delete unused files, move test utilities, update imports, and commit new source files.

**Tech Stack:** Git, bash

---

### Task 1: Delete design iteration files from docs/

**Files:**
- Delete: `docs/dashboard-preview-v1.html`
- Delete: `docs/dashboard-preview-v2.html`
- Delete: `docs/dashboard-preview-v3.html`
- Delete: `docs/login-page-preview.html`
- Delete: `docs/login-page-preview-v2.html`
- Delete: `docs/login-page-preview-v3.html`
- Delete: `docs/logo-badge-dark.svg`
- Delete: `docs/logo-badge-light.svg`
- Delete: `docs/logo-concept-1-island-home.svg`
- Delete: `docs/logo-concept-2-warm-circle.svg`
- Delete: `docs/logo-concept-3-minimal.svg`
- Delete: `docs/logo-concept-4-badge.svg`
- Delete: `docs/logo-concept-5-modern-stack.svg`
- Delete: `docs/warmisle-prototype.html`

- [ ] **Step 1: Delete the files**

```bash
rm docs/dashboard-preview-v1.html \
   docs/dashboard-preview-v2.html \
   docs/dashboard-preview-v3.html \
   docs/login-page-preview.html \
   docs/login-page-preview-v2.html \
   docs/login-page-preview-v3.html \
   docs/logo-badge-dark.svg \
   docs/logo-badge-light.svg \
   docs/logo-concept-1-island-home.svg \
   docs/logo-concept-2-warm-circle.svg \
   docs/logo-concept-3-minimal.svg \
   docs/logo-concept-4-badge.svg \
   docs/logo-concept-5-modern-stack.svg \
   docs/warmisle-prototype.html
```

- [ ] **Step 2: Verify docs/ only contains intended files**

```bash
ls docs/
```

Expected: `prd.md` and `superpowers/` directory only.

- [ ] **Step 3: Commit**

```bash
git add docs/
git commit -m "chore: remove design iteration HTML and logo concept files from docs/"
```

---

### Task 2: Delete redundant build script

**Files:**
- Delete: `build.ps1`

- [ ] **Step 1: Delete the file**

```bash
rm build.ps1
```

- [ ] **Step 2: Commit**

```bash
git add build.ps1
git commit -m "chore: remove redundant PowerShell build script"
```

---

### Task 3: Delete Vite placeholder icon

**Files:**
- Delete: `frontend/public/vite.svg`

- [ ] **Step 1: Delete the file**

```bash
rm frontend/public/vite.svg
```

- [ ] **Step 2: Commit**

```bash
git add frontend/public/vite.svg
git commit -m "chore: remove default Vite placeholder icon"
```

---

### Task 4: Merge frontend test utilities

**Files:**
- Move: `frontend/src/test-utils/auth-mock.ts` → `frontend/src/test/auth-mock.ts`
- Delete: `frontend/src/test-utils/` (empty directory)
- Modify: `frontend/src/views/category/__tests__/Index.test.ts:47`
- Modify: `frontend/src/views/forum/__tests__/Index.test.ts:57`
- Modify: `frontend/src/views/ledger/__tests__/Index.test.ts:54`
- Modify: `frontend/src/views/member/__tests__/Index.test.ts:56`
- Modify: `frontend/src/views/todo/__tests__/Index.test.ts:48`
- Modify: `frontend/src/views/wish/__tests__/Index.test.ts:39`

- [ ] **Step 1: Move the mock file**

```bash
mv frontend/src/test-utils/auth-mock.ts frontend/src/test/auth-mock.ts
rmdir frontend/src/test-utils
```

- [ ] **Step 2: Update imports in 6 test files**

Replace `import '@/test-utils/auth-mock'` with `import '@/test/auth-mock'` in each file:

- `frontend/src/views/category/__tests__/Index.test.ts:47`
- `frontend/src/views/forum/__tests__/Index.test.ts:57`
- `frontend/src/views/ledger/__tests__/Index.test.ts:54`
- `frontend/src/views/member/__tests__/Index.test.ts:56`
- `frontend/src/views/todo/__tests__/Index.test.ts:48`
- `frontend/src/views/wish/__tests__/Index.test.ts:39`

- [ ] **Step 3: Run frontend tests to verify**

```bash
cd frontend && npm test
```

Expected: All tests pass.

- [ ] **Step 4: Commit**

```bash
git add frontend/src/test/ frontend/src/test-utils/ frontend/src/views/
git commit -m "chore: move auth-mock to src/test/, update imports in 6 test files"
```

---

### Task 5: Commit new source files

**Files:**
- Add: `backend/internal/model/local_time.go`
- Add: `frontend/src/components/LogoIcon.vue`
- Add: `frontend/public/favicon.svg`

- [ ] **Step 1: Stage source files**

```bash
git add backend/internal/model/local_time.go \
        frontend/src/components/LogoIcon.vue \
        frontend/public/favicon.svg
```

- [ ] **Step 2: Commit**

```bash
git commit -m "feat: add LocalTime model, LogoIcon component, and favicon"
```

---

### Task 6: Update .gitignore

**Files:**
- Modify: `.gitignore`

- [ ] **Step 1: Append IDE and OS rules**

Add to end of `.gitignore`:

```
# IDE
.idea/
.vscode/

# OS
.DS_Store
Thumbs.db
```

- [ ] **Step 2: Commit**

```bash
git add .gitignore
git commit -m "chore: add IDE and OS patterns to .gitignore"
```

---

### Task 7: Final verification

- [ ] **Step 1: Check Go build**

```bash
cd backend && go build ./...
```

Expected: No errors.

- [ ] **Step 2: Run frontend tests**

```bash
cd frontend && npm test
```

Expected: All tests pass.

- [ ] **Step 3: Verify git status is clean**

```bash
git status
```

Expected: No untracked files (except those intentionally ignored), no unstaged changes.
