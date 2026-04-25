# Common Native Replacements

| Library | Native Alternative | Since |
|---------|-------------------|-------|
| moment.js | Intl.DateTimeFormat + Temporal | ES2020+ |
| lodash.get | Optional chaining (?.) | ES2020 |
| lodash.cloneDeep | structuredClone() | 2022 |
| uuid | crypto.randomUUID() | Node 19+ |
| node-fetch | global fetch() | Node 18+ |
| chalk | util.styleText() | Node 21+ |
| classnames | Template literals | Always |
| underscore | Array methods | ES2015+ |
