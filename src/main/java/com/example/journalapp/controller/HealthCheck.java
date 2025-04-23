package com.example.journalapp.controller;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import java.util.Map;

@RestController
public class HealthCheck{
    @GetMapping("/health-check")
    public String getHealth(){
        return "OK";
    }
    @GetMapping("/__debug/default")
    public Map<String, String> debugDefault(@RequestParam Map<String, String> queryParams) {
        // Return a simple JSON response
        return Map.of("message", "pong", "receivedParams", queryParams.toString());
    }

    @GetMapping("/__debug/optional")
    public Map<String, String> debugOptional(@RequestParam Map<String, String> queryParams) {
        // Return a simple JSON response showing received query parameters
        return Map.of("message", "pong", "receivedParams", queryParams.toString());
    }

    @GetMapping("/__debug/qs")
    public Map<String, String> debugQS(@RequestParam Map<String, String> queryParams) {
        // Return a simple JSON response showing received query parameters
        return Map.of("message", "pong", "receivedParams", queryParams.toString());
    }
}