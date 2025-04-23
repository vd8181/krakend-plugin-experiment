package com.example.journalapp.entity;
//import jakarta.validation.constraints.NotBlank;
import com.fasterxml.jackson.annotation.JsonFormat;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import lombok.Data;
import lombok.Getter;
import lombok.Setter;
import org.bson.types.ObjectId;
import org.springframework.data.annotation.Id;
import org.springframework.data.mongodb.core.index.Indexed;
import org.springframework.data.mongodb.core.mapping.DBRef;
import org.springframework.data.mongodb.core.mapping.Document;
import java.time.LocalDateTime;
import java.util.Date;

@Document(collection = "journal_entries")
@Data
public class JournalEntry {

//    @NotBlank
    @Id
    private ObjectId id;

    @NotBlank
    @Indexed(unique = true)
    private String title;
    @NotBlank
    private String content;
//    @NotNull
//@JsonFormat(shape = JsonFormat.Shape.STRING, pattern = "dd-MM-yyyy")

private Date date;

}
