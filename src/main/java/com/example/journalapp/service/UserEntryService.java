package com.example.journalapp.service;

import com.example.journalapp.entity.JournalEntry;
import com.example.journalapp.entity.UserEntry;
import com.example.journalapp.repository.UserEntryRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;
import org.springframework.stereotype.Service;

import java.util.List;

@Component
@Service
public class UserEntryService  {
    @Autowired
    private UserEntryRepository userEntryRepository;
    public Integer createEntry(UserEntry userEntry){
        try{

            userEntryRepository.insert(userEntry);
            return 1;
        }
        catch(Exception e){
            return 0;
        }
//        return null;
    }
    public List<UserEntry> getAllUserEntries(){
        return userEntryRepository.findAll();
    }
    public Integer updateUsername(String oldUserName, String newUserName){
        UserEntry userEntry = userEntryRepository.findByUsername(oldUserName);
        if(userEntry==null){
            return 0;
        }
        else{
            userEntry.setUsername(newUserName);
            return 1;
        }
    }


    public UserEntry findByUsername(String username) {
        return userEntryRepository.findByUsername(username);
    }

//    public Optional<UserEntry> findById(Long userId) {
//        return userEntryRepository.findById(userId);
//    }

    public Integer deleteEntryByUsername(String userName) {
        try{
            UserEntry userEntry = userEntryRepository.findByUsername(userName);
            if(userEntry==null){
                return 0;
            }
            else{
                userEntryRepository.deleteByUsername(userEntry.getUsername());
                return 1;
            }
        } catch (Exception e) {
            return 0;
        }

    }

    public Integer updateUserEntry(UserEntry user) {
        try{

            userEntryRepository.save(user);
            return 1;
        } catch (Exception e) {
            return 0;
        }
    }


}
