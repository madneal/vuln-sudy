# Secret Scanning 集成方案

本仓库采用“GitHub 原生能力 + 本地/PR 补充检查”的双层方案。

## 1. 能力边界

- GitHub Secret Scanning / Push Protection：由仓库 Settings 开启，在 GitHub 服务端扫描 push 和仓库历史；命中后可以在 push 阶段拒绝更新远端 ref。
- Gitleaks：通过 `.github/workflows/secret-scanning.yml` 在 PR 和 `master` push 上运行，作为可见的 PR 状态检查；通过 `.pre-commit-config.yaml` 提供本地提交前检查。
- `.gitleaks.toml`：在默认规则之外增加本组织测试格式 `TEST_SECRET_[A-Z0-9]{32}`。

Gitleaks 的结果不能替代 GitHub 原生 Secret Scanning；两者规则库、有效性校验、告警生命周期和 Push Protection 行为不同。

## 2. GitHub 侧配置步骤

1. 进入仓库 `Settings → Advanced Security`。
2. 开启 `Secret Protection` / `Secret Scanning`。
3. 开启 `Push Protection`；对测试用例选择 `used in tests`，不要上传真实凭据。
4. 在 `Secret scanning → Custom patterns` 创建组织内部 Secret 格式，例如：

   ```regex
   TEST_SECRET_[A-Z0-9]{32}
   ```

5. 进入 `Settings → Rules → Rulesets`，针对 `master`：
   - 要求 PR 合并；
   - 要求 `Gitleaks`、`CodeQL`、构建和测试检查通过；
   - 对 CodeQL 配置 High 或更高安全告警阻断；
   - 限制绕过权限到安全管理员或指定团队。

## 3. 本地使用

安装 pre-commit 和 Gitleaks 后：

```bash
brew install gitleaks
pipx install pre-commit
pre-commit install
pre-commit run --all-files
```

直接扫描当前工作树：

```bash
gitleaks dir --source . --redact
```

扫描当前分支的 Git 历史：

```bash
gitleaks git --log-opts="origin/master..HEAD" --redact
```

本地测试自定义规则时，可使用以下明确的合成值：

```text
TEST_SECRET_0123456789ABCDEF0123456789ABCDEF
```

该值不对应任何服务，不应被当作凭据使用。

## 4. PR 卡点行为

| 场景 | 本地 hook | Gitleaks PR check | GitHub Push Protection |
|---|---:|---:|---:|
| 新增默认规则命中 | 阻止 commit | 失败 | 可能在 push 时阻断 |
| 新增自定义规则命中 | 阻止 commit | 失败 | 取决于原生 custom pattern 是否已启用 |
| 仅修改已有告警 | 可能再次提示 | 失败或产生重复结果 | 可能不重复阻断 |
| 真实凭据泄露 | 阻止并人工处理 | 失败 | 应立即撤销/轮换 |

## 5. 误报与例外

- 测试假值：关闭为 `used_in_tests`，并在告警中记录测试用途。
- 误报：关闭为 `false_positive`，记录判定依据。
- 真实凭据：先撤销/轮换，再删除代码和 Git 历史；不要用 `false_positive` 代替修复。
- 不建议通过 `paths-ignore` 全局忽略测试目录，除非能保证其中永远不会出现真实凭据。

## 6. 验证清单

1. 在本地运行 `pre-commit run --all-files`，确认发现测试格式。
2. 创建只包含合成值的 PR，确认 `Gitleaks` check 失败并给出文件位置。
3. 向启用 Push Protection 的测试分支执行 `git push`，确认 GitHub 服务端返回 `GH013`。
4. 通过测试用途流程完成一次受控 bypass，确认审计记录存在。
5. 移除测试值并重新提交，确认本地检查、Gitleaks 和 PR 门禁恢复通过。

