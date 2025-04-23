package com.example.journalapp.controller;

//import com.example.journalapp.entity.UserEntry;
import com.example.journalapp.entity.JournalEntry;
import com.example.journalapp.entity.UserEntry;
        import com.example.journalapp.service.UserEntryService;
import jakarta.validation.Valid;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.HashMap;
import java.util.List;

@RestController
@RequestMapping("/user")
public class UserEntryController {
    @Autowired
    private UserEntryService userEntryService;
//    @GetMapping(}")
    @PostMapping("/journal")
    private ResponseEntity createJournalForUser(String username, JournalEntry journalEntry){
        UserEntry userEntry= userEntryService.findByUsername(username);
        return new ResponseEntity(HttpStatus.CREATED);
    }
    @GetMapping()
    private ResponseEntity getAllUsers(){
//  return "Get all users";
//    return  userEntryService.getAllUserEntries();
        try{

            List<UserEntry> userEntries = userEntryService.getAllUserEntries();
            if(userEntries==null){
                return ResponseEntity.status(HttpStatus.NOT_FOUND).body(userEntries);
            }
            else{
                return ResponseEntity.ok(userEntries);

            }
        }
        catch(Exception e){
            return new ResponseEntity<>("Unkown Error Occured.",HttpStatus.INTERNAL_SERVER_ERROR);
//            return new ResponseEntity("Unknown Error occured. Check your connection or it is an error from our side");
        }

    }
    @PostMapping()
    private ResponseEntity createUser(@Valid @RequestBody UserEntry userEntryBody) {
        try {
            Integer response = userEntryService.createEntry(userEntryBody);
            if(response==1){
                return ResponseEntity.ok("User created");
            }
            else{
                return ResponseEntity.status(HttpStatus.BAD_REQUEST).body("User with given username already exists.");
            }
        } catch (Exception e){
            return new ResponseEntity("Unkown error occured while creating the user",HttpStatus.CREATED);
    }
    }
    @PutMapping()
    private ResponseEntity updateUser(@RequestBody HashMap<String,HashMap<String,String>> data){
        try{
            HashMap<String,String> oldData = data.get("oldData");

            HashMap<String,String> newData = data.get("newData");
            //username can never be null in oldData.
            UserEntry user = userEntryService.findByUsername(oldData.get("username"));
            if(newData.get("username")!=null){
                user.setUsername(newData.get("username"));
            }
            if(newData.get("password")!=null){
                user.setPassword(newData.get("password"));
            }
            userEntryService.updateUserEntry(user);
            return ResponseEntity.ok("Update Done");

        }
        catch (Exception e){
        return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body("Unkown Error");
        }
    }
    @DeleteMapping()
    private ResponseEntity deleteUser(@RequestBody UserEntry user){
        try{
            userEntryService.deleteEntryByUsername(user.getUsername());
            return ResponseEntity.ok("Deletion Done");
        }
        catch(Exception e){
            return ResponseEntity.status(HttpStatus.BAD_REQUEST).body("Unkown error occured.");
        }
    }
}
