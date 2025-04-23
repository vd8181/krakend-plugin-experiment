package com.example.journalapp.repository;

import com.example.journalapp.entity.UserEntry;
import org.bson.types.ObjectId;
import org.springframework.data.mongodb.repository.MongoRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface UserEntryRepository extends MongoRepository<UserEntry, ObjectId> {
    UserEntry findByUsername(String username);
    void deleteByUsername(String username);
}
