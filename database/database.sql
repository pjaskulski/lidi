CREATE DATABASE dictionary;

USE dictionary;

DROP TABLE engpol;

CREATE TABLE engpol (
	id INT AUTO_INCREMENT PRIMARY KEY,
	english VARCHAR(100),
    polish VARCHAR(100)
);

CREATE INDEX idx_english ON engpol(english); 
CREATE INDEX idx_polish ON engpol(polish);

LOAD DATA INFILE 'words.csv'
INTO TABLE engpol
CHARACTER SET utf8
FIELDS TERMINATED BY ','
OPTIONALLY ENCLOSED BY '"'
LINES TERMINATED BY '\n'
(english, polish);

-- Test:
SELECT polish FROM engpol where english='richer';
SELECT COUNT(*) FROM engpol;

-- eksport danych
SELECT 
    english, polish 
FROM engpol 
INTO OUTFILE 'words_export.csv'
CHARACTER SET utf8
FIELDS TERMINATED BY ','
OPTIONALLY ENCLOSED BY '"'
LINES TERMINATED BY '\n';
