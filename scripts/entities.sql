CREATE TABLE Users (
    Username TEXT NOT NULL,
    Password TEXT NOT NULL,
    Email TEXT NOT NULL,
    UID TEXT NOT NULL,
    Occupation TEXT NOT NULL,
    CreatedAt TIMESTAMP DEFAULT current_timestamp,
    DeletedAt TIMESTAMP,
    UpdatedAt TIMESTAMP DEFAULT current_timestamp,
    PRIMARY KEY(UID)
);

CREATE TABLE Students (
    Gender TEXT NOT NULL, 
    DOB TEXT NOT NULL,
    SchoolYear INTEGER NOT NULL,
    CurrentLessonID TEXT NOT NULL,
    UID TEXT NOT NULL,
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

CREATE TABLE Units (
    UnitName TEXT NOT NULL,
    UnitID TEXT,
    UnitDescription TEXT,
    CreatedAt TIMESTAMP DEFAULT current_timestamp,
    UpdatedAt TIMESTAMP,
    DeletedAt TIMESTAMP DEFAULT current_timestamp,
    PRIMARY KEY(UnitID)
);

CREATE TABLE Chapters (
    ChapterID TEXT,
    ChapterName TEXT,
    ChapterDescription TEXT,
    UnitID Text NOT NULL,
    FOREIGN KEY(UnitID) REFERENCES Units,
    PRIMARY KEY(ChapterID)
);

CREATE TABLE Lessons (
    LessonName TEXT,
    LessonContent TEXT[][],
    LessonID TEXT,
    ChapterID TEXT,
    CreatedAt TIMESTAMP DEFAULT current_timestamp,
    DeletedAt TIMESTAMP,
    UpdatedAt TIMESTAMP DEFAULT current_timestamp,
    MinimumScoreToPass INTEGER,
    FOREIGN KEY(ChapterID) REFERENCES Chapters,
    PRIMARY KEY(LessonID)
);

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

CREATE TABLE LessonsCompleted (
    LessonID TEXT NOT NULL,
    UID TEXT NOT NULL,
    FOREIGN KEY(LessonID) REFERENCES Lessons,
    FOREIGN KEY(UID) REFERENCES Users,
    PRIMARY KEY(LessonID, UID)
);

CREATE TABLE ChaptersCompleted (
    ChapterID TEXT NOT NULL,
    UID TEXT NOT NULL,
    FOREIGN KEY(ChapterID) REFERENCES Chapters,
    FOREIGN KEY(UID) REFERENCES Users,
    PRIMARY KEY(UID, ChapterID)
);

CREATE TABLE UnitsCompleted (
    UnitID TEXT NOT NULL,
    UID TEXT NOT NULL,
    FOREIGN KEY(UnitID) REFERENCES Units,
    FOREIGN KEY(UID) REFERENCES Users,
    PRIMARY KEY(UnitID, UID)
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