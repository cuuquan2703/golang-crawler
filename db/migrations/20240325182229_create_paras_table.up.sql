create table if not exists Paragraphs (
    id int not null,
    url varchar(255),
    json varchar(255),
    linecount int,
    wordcount int,
    charcount int,
    avglength int,
    wordfreq varchar(100000)
)