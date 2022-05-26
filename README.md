## Import procurement

When a procurement is imported successfully, import status is given to this file "pending", "imported", "failed".

**pending**

The worker will routinely import pending files to DB.


**imported**

The worker will change the file status to imported once file data has been imported to DB successfully.

**failed**

The worker will change the file status to failed if file data failed to import to DB. The reason will be stated at the column `failed_reason` .


# TODOs

- [x] auth middleware to validate jwt token.
- [x] reserve available stock for user.
- [] move 'imported' inventory to another bucket.
