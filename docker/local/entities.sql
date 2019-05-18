CREATE TABLE Users (
    Username TEXT NOT NULL,
    Password TEXT NOT NULL,
    Email TEXT NOT NULL,
    UID TEXT NOT NULL,
    Occupation TEXT NOT NULL,
    CreatedAt TIMESTAMP DEFAULT current_timestamp,
    DeletedAt TIMESTAMP,
    UpdatedAt TIMESTAMP DEFAULT current_timestamp,
    FirstName TEXT NOT NULL, 
    LastName TEXT NOT NULL,
    PRIMARY KEY(UID)
);

CREATE TABLE GameText (
    ID TEXT NOT NULL, 
    Txt TEXT NOT NULL,
    PRIMARY KEY(ID)
);

-- INSERT INTO Users(
--     Username,
--     Password,
--     Email,
--     UID,
--     Occupation,
--     FirstName,
--     LastName
-- ) VALUES (
--     'tester',
--     'password',
--     'tester@gmail.com',
--     '1',
--     'Student',
--     'Tester',
--     'McTestem'
-- );

-- INSERT INTO Users(
--   Username,
--   Password,
--   Email,
--   UID,
--   Occupation,
--   FirstName,
--   LastName
-- ) VALUES (
--            'tester',
--            'password',
--            'tester@gmail.com',
--            '2',
--            'Student',
--            'Tester',
--            'McTestem'
-- );

CREATE TABLE Students (
    Gender TEXT NOT NULL, 
    DOB TEXT NOT NULL,
    CurrentLessonID TEXT NOT NULL,
    CurrentChapterName TEXT NOT NULL, 
    CurrentChapterID TEXT NOT NULL,
    UID TEXT NOT NULL,
    FOREIGN KEY(UID) REFERENCES Users
);

CREATE TABLE Pupils (
    SchoolYear TEXT,
    UID TEXT,
    FOREIGN KEY(UID) REFERENCES Users
);

CREATE TABLE Instructors (
    Gender TEXT NOT NULL,
    DOB TEXT NOT NULL,
    SchoolType TEXT NOT NULL,
    SchoolName TEXT NOT NULL,
    UID TEXT NOT NULL,
    FOREIGN KEY(UID) REFERENCES Users
);

CREATE TABLE Employed (
    Occupation TEXT,
    UID TEXT NOT NULL,
    Primary Key(UID)
);


CREATE TABLE Chapters (
    ChapterID TEXT,
    ChapterName TEXT,
    ChapterDescription TEXT,
    ChapterImage Text,
    NextChapterID TEXT,
    PRIMARY KEY(ChapterID)
);

-- INSERT INTO Chapters (
--     ChapterID
-- ) VALUES (
--     '1'
-- );

CREATE TABLE Lessons (
    LessonName TEXT,
    LessonText TEXT[],
    LessonID TEXT,
    ChapterID TEXT,
    Image Text[],
    LessonDescriptions TEXT[],
    CreatedAt TIMESTAMP DEFAULT current_timestamp,
    DeletedAt TIMESTAMP,
    UpdatedAt TIMESTAMP DEFAULT current_timestamp,
    MinimumScoreToPass INTEGER[],

    NextLessonID TEXT,

    FOREIGN KEY(ChapterID) REFERENCES Chapters,
    PRIMARY KEY(LessonID)
);

-- INSERT INTO Lessons (
--     LessonID,
--     ChapterID
-- ) VALUES (
--     '2',
--     '1'
-- );

CREATE TABLE Schools (
    SchoolName TEXT NOT NULL,
    SchoolID TEXT,
    PRIMARY KEY(SchoolID)
);

CREATE TABLE Classrooms (
    ClassroomID TEXT,
    InstructorID TEXT NOT NULL,
    Year INTEGER NOT NULL,
    Subject TEXT NOT NULL,
    SchoolID TEXT NOT NULL,
    FOREIGN KEY(InstructorID) REFERENCES Users,
    FOREIGN KEY(SchoolID) REFERENCES Schools,
    PRIMARY KEY(ClassroomID)
);

CREATE TABLE lessonscompleted (
    LessonID TEXT NOT NULL,
    UID TEXT NOT NULL,
    Accuracy DECIMAL,
    WPM DECIMAL,
    ChapterID TEXT,
    FOREIGN KEY(LessonID) REFERENCES Lessons,
    FOREIGN KEY(UID) REFERENCES Users,
    PRIMARY KEY(LessonID, UID)
);

-- INSERT INTO lessonscompleted(LessonID, UID, ChapterID, Accuracy, WPM) VALUES ('2', '1', '1', '100', '100');

CREATE TABLE ChaptersCompleted (
    ChapterID TEXT NOT NULL,
    UID TEXT NOT NULL,
    FOREIGN KEY(ChapterID) REFERENCES Chapters,
    FOREIGN KEY(UID) REFERENCES Users,
    PRIMARY KEY(UID, ChapterID)
);

CREATE TABLE Enrolled (
    StudentID TEXT NOT NULL,
    InstructorID TEXT NOT NULL,
    SchoolID TEXT NOT NULL,
    ClassroomID TEXT NOT NULL,
    FOREIGN KEY(StudentID) REFERENCES Users,
    FOREIGN KEY(InstructorID) REFERENCES Users,
    FOREIGN KEY(SchoolID) REFERENCES Schools,
    FOREIGN KEY(ClassroomID) REFERENCES Classrooms,
    PRIMARY KEY(StudentID, ClassroomID)
);

CREATE TABLE ProgressReport (
    StudentID TEXT NOT NULL,
    WPM INTEGER,
    Accuracy INTEGER,
    FOREIGN KEY(StudentID) REFERENCES Users,
    PRIMARY KEY(StudentID)
);