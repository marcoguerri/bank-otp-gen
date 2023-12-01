Code which generates 2FA OTP for an Italian bank, provided the necessary
secrets are known. This code was reverse engineered from the bank's Android 
application for disaster recovery purposes. 
See the corresponding [blog post](https://marcoguerri.github.io/2023/09/09/android-home-banking.html) 
for full details on the reverse engineering exercise.

# Usage

As outlined in the [blog post](https://marcoguerri.github.io/2023/09/09/android-home-banking.html), OTP calculation requires
a pin and two keys, `sc_sac` and `sc_k2`. Given a symmetrically encrypted file (`secrets.gpg`) with the following format:

```
{
  "pin": "00000",
  "sc_sac": "<KEY>",
  "sc_k2": "<KEY>"
}
```

OTP can be obtained with the following command:

```
gpg -d secrets.gpg | go run main.go -qrcode <QRCODE_DATA> -minTimeDelta -2 -maxTimeDelta  2
```
where `<QRCODE_DATA>` has been stripped off the `SC:` prefix. The command above makes the assumption that `<QRCODE_DATA>` has
been generated between `now()-2h` and `now()+2h`.
