 operator-sdk init --domain clodevo.com --repo github.com/lakhlaifi/configprobe-operator


 operator-sdk create api --group synthetic --version v1 --kind ConfigProbe --resource --controller


 make generate
