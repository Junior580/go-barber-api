const crypto = require("crypto");

// 1️⃣ Gera AES + IV
const aesKey = crypto.randomBytes(32); // AES-256
const iv = crypto.randomBytes(16); // nonceSize = 16

const body = {
  action: "INIT",
};

// 2️⃣ Public Key
const publicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAwVs4eSGLNfuikQdJG5Gu
rh4DsdJqg2yypPVfnJkgBPf+7dyk8vlnYU/88MJc1FCQ+KH7mAeeMgjP6zuY2bn3
NLXSvJR9xEec5nojsYCboJVVlUjIXXDPqZoWAwJ5IOrVuDimXDUI/aQQkTQpDKvh
L4Mm4OI6rINqKX9albdHz224EiKH3+npjhyIo4a8HpS6Sql5rq/6KcFGCniR710B
0ZP60kG0qSo8kG7cej1zMN+qYVzISftSMCSgO8Yd0bZSLoI28qTrDmKuo7dFoKm+
IDzBMIU7+TvdhBXU2V/Dr1hwKUxo3DpuaDkpYZSaHEFnMf6syY2QcRrLuKhg3kiX
kwIDAQAB
-----END PUBLIC KEY-----`;

// 3️⃣ AES-256-GCM encrypt
const cipher = crypto.createCipheriv("aes-256-gcm", aesKey, iv);

const encrypted = Buffer.concat([
  cipher.update(JSON.stringify(body), "utf8"),
  cipher.final(),
]);

const authTag = cipher.getAuthTag();

// ciphertext + tag (igual no Go)
const encryptedFlowData = Buffer.concat([encrypted, authTag]);

// 4️⃣ RSA-OAEP-SHA256 encrypt da AES key
const encryptedAESKey = crypto.publicEncrypt(
  {
    key: publicKey,
    padding: crypto.constants.RSA_PKCS1_OAEP_PADDING,
    oaepHash: "sha256",
  },
  aesKey,
);

// 5️⃣ Payload FINAL (BASE64)
const payload = {
  encrypted_aes_key: encryptedAESKey.toString("base64"),
  encrypted_flow_data: encryptedFlowData.toString("base64"),
  initial_vector: iv.toString("base64"),
};

console.log(JSON.stringify(payload, null, 2));
