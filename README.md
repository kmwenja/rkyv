rkyv
====

`rkyv` archive file format and binary that manages these files.

File Format
-----------

- meta pointer
- file bytes
- meta:
  - version
  - uuid
  - date_created
  - date_updated
  - files:
    - name
    - type
    - size
    - hash
  - tags
  - search

Cli
---

- rkyv create -i -t 'tag1,tag2,tag3' -s 'search1,search2,search3' file1 file2
- rkyv update -i -t -s <rkyv_file> file1 file2
- rkyv extract <rkyv_file> <file index>
- rkyv info <rkyv_file>
- rkyv list -t -T -s -S -p -c -b -a
- rkyv scan -f <index_file>
- rkyv index TODO
