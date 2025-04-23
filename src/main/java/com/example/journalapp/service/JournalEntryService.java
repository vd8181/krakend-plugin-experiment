package com.example.journalapp.service;

import com.example.journalapp.entity.JournalEntry;
import com.example.journalapp.repository.JournalEntryRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;
import org.springframework.stereotype.Service;

import java.time.LocalDate;
import java.util.Date;
import java.util.List;
import java.util.Optional;

//@Component
@Service
public class JournalEntryService  {
    @Autowired
    private JournalEntryRepository journalEntryRepository;
    public void createEntry(JournalEntry journalEntry){
        journalEntryRepository.insert(journalEntry);
    }
    public List<JournalEntry> getJournalEntries(){
        return journalEntryRepository.findAll();
    }
    public Integer updateEntry(JournalEntry oldEntry , JournalEntry newEntry){
//         journalEntryRepository.save();
        JournalEntry oldJournalEntryToBeUpdated = journalEntryRepository.findByTitle(oldEntry.getTitle());
        if(oldJournalEntryToBeUpdated!=null){
            oldJournalEntryToBeUpdated.setDate(new Date());
            oldJournalEntryToBeUpdated.setContent(newEntry.getContent());
            oldJournalEntryToBeUpdated.setTitle(newEntry.getTitle());
            journalEntryRepository.save(oldJournalEntryToBeUpdated);
            return 1;
        }
        else{
            return 0;
        }
    }


    public JournalEntry findByTitle(String title) {
        return journalEntryRepository.findByTitle(title);
    }

//    public Optional<JournalEntry> findById(Long journalId) {
//        return journalEntryRepository.findById(journalId);
//    }

    public Optional<String> deleteEntryByTitle(String journalTitle) {
        JournalEntry check = journalEntryRepository.findByTitle(journalTitle);
        if(check!=null){

            journalEntryRepository.deleteByTitle(journalTitle);
            return Optional.ofNullable(journalTitle);
        }
        else{
            return Optional.empty();
        }

//        return ;
    }
}
