package com.example.journalapp.service;

import com.example.journalapp.entity.DataDTO;
import java.security.PublicKey;
import java.security.Signature;
import java.util.Base64;

public class CertificateValidationService {

    private PublicKey serverPublicKey;

    // Constructor to load the public key from a file
    public CertificateValidationService(String publicKeyPath) {
        try {
            // Load the public key dynamically
            this.serverPublicKey = KeyLoader.loadPublicKey(publicKeyPath);
        } catch (Exception e) {
            throw new RuntimeException("Failed to load public key", e);
        }
    }

    /**
     * Validates the certificate (signature) using the server's public key and checks if it matches the payload.
     * @param data The nested data field to be validated.
     * @param signature The signature to be decoded and verified.
     * @return HTTP status code: 200 if valid, 400 if invalid.
     */
    public int validateCertificateFunction(DataDTO data, String signature) {
        try {
            // Convert the DataDTO object to a string representation of the payload
            String payload = convertDataDTOToPayload(data);

            // Decode the signature from Base64
            byte[] decodedSignature = Base64.getDecoder().decode(signature);

            // Initialize the Signature object with the server's public key
            Signature sig = Signature.getInstance("SHA256withRSA");
            sig.initVerify(serverPublicKey);

            // Provide the serialized payload to the Signature object for verification
            sig.update(payload.getBytes());

            // Verify the signature against the payload
            boolean isValid = sig.verify(decodedSignature);

            // Return HTTP status code based on validation result
            return isValid ? 200 : 400;

        } catch (Exception e) {
            e.printStackTrace();
            // Return 400 in case of any errors
            return 400;
        }
    }

    /**
     * Helper method to convert a DataDTO object to a string representation (e.g., JSON-like format).
     * @param data The DataDTO object.
     * @return A string representation of the payload.
     */
    private String convertDataDTOToPayload(DataDTO data) {
        // Serialize the nested `data` field into JSON
        return String.format("{\"name\":\"%s\",\"email\":\"%s\",\"dob\":\"%s\"}",
                data.getName(), data.getEmail(), data.getDob());
    }
}
