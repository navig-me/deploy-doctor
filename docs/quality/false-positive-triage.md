# False Positive Triage and Default Tuning

## Workflow
1. Reproduce finding.
2. Confirm expected platform/runtime behavior.
3. Label finding: true positive / false positive / ambiguous.
4. For false positives:
   - tighten heuristic
   - lower default severity when confidence is low
   - add regression fixture test
5. Track changes in changelog.

## Required Artifacts
- before/after sample output
- test covering tuned behavior
- rationale for severity/default change
