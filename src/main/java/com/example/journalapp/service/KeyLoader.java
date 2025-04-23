package com.example.journalapp.service;

import java.io.File;
import java.io.FileInputStream;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.security.KeyFactory;
import java.security.PublicKey;
import java.security.spec.X509EncodedKeySpec;
import java.util.Base64;

public class KeyLoader {

    /**
     * Loads a public key from a file and converts it into a PublicKey object.
     * Assumes the public key file is Base64 encoded (e.g., PEM format).
     * @param filePath Path to the public key file.
     * @return The PublicKey object.
     * @throws Exception If there's an error loading or parsing the key.
     */
    public static PublicKey loadPublicKey(String filePath) throws Exception {
        // Read all bytes from the public key file
        byte[] keyBytes = Files.readAllBytes(Paths.get(filePath));

        // Convert the key bytes into a string and remove any PEM header/footer
        String keyString = new String(keyBytes).replace("-----BEGIN PUBLIC KEY-----", "")
                .replace("-----END PUBLIC KEY-----", "")
                .replaceAll("\\s", ""); // Remove whitespace

        // Decode the Base64-encoded string into raw key bytes
        byte[] decodedKey = Base64.getDecoder().decode(keyString);

        // Generate the PublicKey object from the decoded bytes
        X509EncodedKeySpec spec = new X509EncodedKeySpec(decodedKey);
        KeyFactory keyFactory = KeyFactory.getInstance("RSA");
        return keyFactory.generatePublic(spec);
    }
}
