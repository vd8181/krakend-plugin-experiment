package com.example.journalapp.controller;

import com.example.journalapp.entity.JournalEntry;
import com.example.journalapp.entity.UserEntry;
import com.example.journalapp.service.JournalEntryService;
//import jakarta.validation.Valid;
import com.example.journalapp.service.UserEntryService;
import jakarta.validation.Valid;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.Date;
import java.util.List;
import java.util.Optional;

//import javax.validation.Valid;
//import ja
@RestController //Rest Controller offers Component + REST Functionality
@RequestMapping("/journals")
public class JournalEntityController {
//    SimpleDateFormat dateFormat = new SimpleDateFormat("yyyy-MM-dd HH:mm:ss");

    @Autowired
    private JournalEntryService journalEntryService;
    @Autowired
    private UserEntryService userEntryService;
    @GetMapping
    public List<JournalEntry> getAllJournals(){

        System.out.println("IN GET ALL JOURNALS ");
        return journalEntryService.getJournalEntries();
    }
    @PostMapping
    public ResponseEntity<?> createEntry(@Valid @RequestBody JournalEntry journalEntry){
        try{
            JournalEntry duplicateJournalEntry = journalEntryService.findByTitle(journalEntry.getTitle());
            if(duplicateJournalEntry!=null){
//                return "Journal with title already exists.";
//                return new ResponseEntity<>("Journal with given title already exists.", HttpStatus.BAD_REQUEST);
                    return ResponseEntity.status(HttpStatus.BAD_REQUEST).body("Journal with given title already exists.");
            }
            Date currentDate = new Date();
            journalEntry.setDate(currentDate);
            String formattedDate = currentDate.toString();

            journalEntryService.createEntry(journalEntry);
//            return  String.format("Journal entry %s %s %s Created",journalEntry.getTitle(),journalEntry.getContent(),formattedDate);
            return ResponseEntity.ok(journalEntry);
        }
        catch(Exception e){
//            return "ID already exists";
//        return new ResponseEntity<>(e.toString(), HttpStatus.BAD_REQUEST);
            return ResponseEntity.status(HttpStatus.BAD_REQUEST).body("Journal ID already exists.");
        }
    }
    @PutMapping()
    public  ResponseEntity<?> updateEntry(@RequestBody JournalEntry oldEntry , @RequestBody JournalEntry newEntry){
//        JournalEntry oldVersion = journalEntryService.findById(oldEntry.getId());
//        journalEntryService.updateEntry(oldEntry,newEntry);
//        return "Updated";
        try{

            Integer ans = journalEntryService.updateEntry(oldEntry,newEntry);
            if(ans==1){

                return ResponseEntity.ok(newEntry);
            }
            else{
                return ResponseEntity.status(HttpStatus.BAD_REQUEST).body("Invalid oldEntry");
            }
        }
        catch(Exception e){
//            return new ResponseEntity<>("Unknown error occured"+e.toString(),HttpStatus.BAD_REQUEST);
                return ResponseEntity.status(HttpStatus.BAD_REQUEST).body("Unkown error occured.");
        }
        //Dont have to update t

    }

    @GetMapping("{username}")
    public ResponseEntity<?> getJournalEntryByUsername(@PathVariable("username") String user){
//        System.out.println("HERE");
        try{
            List<JournalEntry> journalEntries = journalEntryService.getJournalEntries();
           return new ResponseEntity<>(journalEntries,HttpStatus.OK);
        }
        catch(Exception e){
    return new ResponseEntity<>(e,HttpStatus.INTERNAL_SERVER_ERROR);
        }
    }
    @DeleteMapping("{title}")
    public ResponseEntity<?> deleteEntry(@PathVariable("title") String journalTitle){
        try{
            Optional<String> journalEntryTitle = journalEntryService.deleteEntryByTitle(journalTitle);
//            System.out.println(journalEntryTitle.orElse(null));
            if (journalEntryTitle.isPresent()) {
                return ResponseEntity.ok(journalTitle);
            } else {
                return ResponseEntity.status(HttpStatus.NOT_FOUND).body("Journal with given title not found.");
            }
        }
        catch(Exception e){
        return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(e);
        }
    }
}
