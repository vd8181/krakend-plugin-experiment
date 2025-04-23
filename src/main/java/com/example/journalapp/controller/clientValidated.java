package com.example.journalapp.controller;

import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/clientValidated")
public class clientValidated {

    @GetMapping
    public ResponseEntity<String> clientValidated() {
        System.out.println("CLIENT VALIDATED");
        return ResponseEntity.status(HttpStatus.OK).body("{\"message\": \"Certificate validation succeeded.\"}");
    }
}