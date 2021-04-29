# lt-text-to-ipa

[![Go](https://github.com/airenas/lt-text-to-ipa/actions/workflows/go.yml/badge.svg)](https://github.com/airenas/lt-text-to-ipa/actions/workflows/go.yml) [![Coverage Status](https://coveralls.io/repos/github/airenas/lt-text-to-ipa/badge.svg?branch=dev)](https://coveralls.io/github/airenas/lt-text-to-ipa?branch=main)

Services for Lithuanian text to IPA transcription. Full API documentation is provided [here](https://app.swaggerhub.com/apis/aireno/lt-text-to-ipa).

## Samples

For a text:

```bash
curl http://klcdocker.vdu.lt/text-transcriber/ipa -X POST -d 'Kairas, ne, vos ne vos.'
```

For one word:
```bash
curl http://klcdocker.vdu.lt/text-transcriber/ipa/Jonai -X GET 
```

## Build 

```bash
make generate
make test
cd build && make clean dbuild dpush
```

## Deploy

### Local

It depends on some private repos and tools: *accenter, transcriber, IPA converter*. Local deployment from zero is not possible.

### KLC

Copy scripts and redeploy:
```bash
cd deploy/klc && make restart
```

---
### Author

**Airenas Vaičiūnas**
 
* [github.com/airenas](https://github.com/airenas/)
* [linkedin.com/in/airenas](https://www.linkedin.com/in/airenas/)


---
### License

Copyright © 2021, [Airenas Vaičiūnas](https://github.com/airenas).
Released under the [The 3-Clause BSD License](LICENSE).

---
