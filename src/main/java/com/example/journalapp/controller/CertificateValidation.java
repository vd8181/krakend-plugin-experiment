//package com.example.journalapp.controller;
//
//import com.example.journalapp.entity.DataDTO;
//import com.example.journalapp.entity.Payload;
//import com.example.journalapp.service.CertificateValidationService;
//import org.springframework.http.HttpStatus;
//import org.springframework.http.ResponseEntity;
//import org.springframework.web.bind.annotation.*;
//
//@RequestMapping("/validatecertificate")
//@RestController
//public class CertificateValidation {
//
//    private CertificateValidationService certificateValidationService;
//
//    // Initialize the service with the path to the public key file
//    public CertificateValidation() {
//        String publicKeyPath = "C:/Users/vd8181/Desktop/journalapp/src/certs/demo2/public_key.pem"; // Update with your actual path
//        this.certificateValidationService = new CertificateValidationService(publicKeyPath);
//    }
//
//    @PostMapping
//    public ResponseEntity<?> validateCertificate(@RequestBody Payload payload) {
//        try {
//            // Extract the nested `data` field and the `signature` field
//            DataDTO data = payload.getData();
//            String signature = payload.getSignature();
//
//            // Validate the payload using the service
//            int validationResult = certificateValidationService.validateCertificateFunction(data, signature);
//
//            // Return 200 status with a JSON response if validation succeeds
//            if (validationResult == 200) {
//                return ResponseEntity.status(HttpStatus.OK).body("{\"message\": \"Certificate validation succeeded.\"}");
//            } else {
//                // Return 400 status with a JSON response if validation fails
//                return ResponseEntity.status(HttpStatus.BAD_REQUEST).body("{\"error\": \"Certificate validation failed.\"}");
//            }
//        } catch (Exception e) {
//            // Log the exception and return a 500 error response
//            e.printStackTrace();
//            return ResponseEntity.status(500).body("{\"error\": \"Internal server error.\"}");
//        }
//    }
//}
