# Surveys Provider

This tiny server provides a list of available json files from a directory tree and expose the list and the file with a simple API. It's dedicated to be used with the [survey-viewer application]() to facilitate an online (users don't need to run the application locally and have the surveys on their device to view and test the surveys).

The service expose the list of available json files in a directory tree (scans the subdirectories).

Limitations:
- No check is currently done on the json files exposed, they can be surveys or not, it's the user's to check the files in the exposed directory are valid and wanted to be exposed.
- Update operation is not safe with concurrent read

# Configuration

The service expects its configuration from environment variables :

- SURVEY_DIR : Path to the root directory where the json surveys files are placed. they can be in subdirectories

# API 

URI are from base URL of the exposed services

- [baseURL]/list : List of the available json files as a JSON array of file descriptor
- [baseURL]/survey?id={id} : loads the survey identified by the provided {id} (id field of the file descriptor)
- [baseURL]/update: hot reload of the available files in the directory (without reloading the service)

