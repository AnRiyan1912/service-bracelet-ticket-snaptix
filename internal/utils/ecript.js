const crypto = require("crypto");

function encryptShort(data, secret) {
  const key = Buffer.from(secret.padEnd(16, "0"), "utf8"); // Pastikan 16-byte key (AES-128)
  const nonce = Buffer.from(secret.slice(0, 6).padEnd(16, "0"), "utf8"); // Nonce dari secret
  const cipher = crypto.createCipheriv("aes-128-ctr", key, nonce);

  const ciphertext = cipher.update(data, "utf8", "hex") + cipher.final("hex");
  const encryptedInt = BigInt("0x" + ciphertext);

  // Base32 encoding lalu ambil 6 karakter
  const shortCode = encryptedInt.toString(36).slice(0, 6).toUpperCase();
  return shortCode;
}

function decryptShort(shortCode, secret) {
  const key = Buffer.from(secret.padEnd(16, "0"), "utf8"); // Pastikan 16-byte key
  const nonce = Buffer.from(secret.slice(0, 6).padEnd(16, "0"), "utf8"); // Nonce dari secret

  // Decode dari Base36 ke hex
  const encryptedInt = BigInt("0x" + parseInt(shortCode, 36).toString(16));
  const encryptedHex = encryptedInt.toString(16).padStart(16, "0");

  const decipher = crypto.createDecipheriv("aes-128-ctr", key, nonce);
  const decryptedText =
    decipher.update(encryptedHex, "hex", "utf8") + decipher.final("utf8");

  return decryptedText;
}
