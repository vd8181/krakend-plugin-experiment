package com.example.journalapp.entity;


import jakarta.validation.constraints.NotBlank;
import lombok.Data;
import lombok.NonNull;
import org.bson.types.ObjectId;
import org.springframework.data.annotation.Id;
import org.springframework.data.mongodb.core.index.Indexed;
import org.springframework.data.mongodb.core.mapping.DBRef;
import org.springframework.data.mongodb.core.mapping.Document;

import java.util.ArrayList;
import java.util.List;

@Document(collection = "users")
@Data
public class UserEntry {
    @Id
    private ObjectId id;
    @NonNull
    @Indexed(unique = true)
    private String username;
//    @NotBlank
@NonNull
    private String password;
    @DBRef
    private List<JournalEntry> journalEntries = new ArrayList<>();

}
