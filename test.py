import os
from pyvalkyrie import Pyvalkyrie
from pyvalkyrie.types import Pagination, JobListExecutionsResponse, JobRetrieveJobExecutionsResponse

client = Pyvalkyrie(
    api_key="abcd"  # This is the default and can be omitted
)

print (client.languages.list())