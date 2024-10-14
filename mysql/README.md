# MySQL Configuration

this directory contains essential files for establishing
the correct database structure.

| File                   | Purpose                                                                                         |
|------------------------|-------------------------------------------------------------------------------------------------|
| my.cnf                 | MySQL Configuration file, same as default, but including extra permissions for LOAD DATA INFILE |
| mytempo.sql            | Build the database schema                                                                       |
| fix\_file\_permissions | Patches mytempo.sql to grant FILE permissions (essentially enabling LOAD DATA INFILE)           |

## Workflow

- run `./fix_file_permissions`
- Done! You can go to `..` and `docker compose up`

### NOTE: remember to clean up an existing build/mysql volume before restarting

