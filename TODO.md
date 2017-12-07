# TODO
- `GetInitialMigration` should returns `(up, down string)` so we don't have to export `Migration.Up` and `Migration.Down`, never
- Define `Begin` and `End` in the `Adapter` interface to provides better transaction boundaries
- Rewrite `run` for better lisibility
- Redefines what really needs to be exported such as `Migration`'fields to keep the outer API simple and developer friendly