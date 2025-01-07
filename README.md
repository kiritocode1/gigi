![](<banner.png>)

## Core Storage
- Object storage:
  - [x] Blob: Store the content of a file `Blob.Serialize(), Blob.Hash()`
  - [x] Tree: Store the content of a directory `Tree.Serialize(), Tree.Hash(), ParseTree()`
  - [x] Commit: Store the content of a commit `Commit.Serialize(), Commit.Hash(), ParseCommit()`
  - [x] Object compression using zlib `Repository.WriteObject(), Repository.ReadObject()`
  - [x] Object validation and size limits `isValidObjectType(), ValidateHash()`
  - [x] Cross-platform file metadata handling `getFileMetadata()` (linux.go/win.go)

## Index Management
- Index:
  - [x] Store file metadata in the index `IndexEntry struct`
  - [x] Store file content hash in the index `writeIndexEntries()`
  - [x] Binary index format with versioning `WriteIndex(), ReadIndexFiles()`
  - [x] Index checksum validation `ReadIndexFiles()`
  - [x] 8-byte alignment padding `writeIndexEntries(), readIndexEntries()`

## Commit Operations
- Commit:
  - [x] Store commit metadata `NewCommit(), Sign struct`
  - [x] Store commit message `Commit.Serialize()`
  - [x] Support multiple parent commits `Commit.ParentHash[]`
  - [x] Store tree hash `Commit.TreeHash`
  - [x] Update HEAD reference `Repository.UpdateHEAD()`
  - [x] Timezone handling in signatures `Sign.String()`

## Tree Operations
- Tree:
  - [x] Store directory structure `Tree.entries`
  - [x] Support multiple file modes `isValidMode()`
  - [x] Sort entries for consistent hashing `Tree.Serialize()`
  - [x] Binary tree format `Tree.Serialize(), ParseTree()`
  - [x] Tree entry validation `ParseTree()`

## Repository Management
- Repository:
  - [x] Initialize new repository `InitRepository()`
  - [x] Write objects with SHA1 hashing `HashObject(), HashFile(), HashContent()`
  - [x] Read objects with type validation `Repository.ReadObject()`
  - [x] Add files to staging `Repository.AddFile()`
  - [x] Commit changes `Repository.Commit()`
  - [x] Update HEAD reference `Repository.UpdateHEAD()`

## Work in Progress
- Remote Operations:
  - [ ] Push to remote repository `Push()` (defined in interface)
  - [ ] Pull from remote repository `Pull()` (defined in interface)
  - [ ] Clone repository `Clone()` (defined in interface)
  - [ ] Fetch updates (not implemented)
  
- Branch Management:
  - [ ] Create branches (not implemented)
  - [ ] Switch branches (not implemented)
  - [ ] Merge branches (not implemented)
  - [ ] Delete branches (not implemented)

- Log and History:
  - [ ] View commit history `Log()` (defined in interface)
  - [ ] Show commit details (not implemented)
  - [ ] Display file changes (not implemented)

## Additional Features Needed
- Working Directory:
  - [ ] Status tracking (not implemented)
  - [ ] Dirty state detection (not implemented)
  - [ ] Ignore patterns (.ggignore) (not implemented)
  
- Recovery and Maintenance:
  - [ ] Garbage collection (not implemented)
  - [ ] Object verification (not implemented)
  - [ ] Repository repair tools (not implemented)

- Configuration:
  - [ ] User settings (not implemented)
  - [ ] Repository config (not implemented)
  - [ ] Global config (not implemented)