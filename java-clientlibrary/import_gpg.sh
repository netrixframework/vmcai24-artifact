echo -n "$GPG_SIGNING_KEY" | base64 --decode | gpg --import
gpg --keyring secring.gpg --export-secret-keys > ~/secring.gpg