package com.example.journalapp.entity;

public class Payload {
    private DataDTO data; // Nested data field
    private String signature; // Signature field

    // Getters and Setters
    public DataDTO getData() {
        return data;
    }

    public void setData(DataDTO data) {
        this.data = data;
    }

    public String getSignature() {
        return signature;
    }

    public void setSignature(String signature) {
        this.signature = signature;
    }
}
