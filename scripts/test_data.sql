INSERT INTO Users(Username, Password, Email, UID) 
VALUES('ibarrac27', 'password', 'ibarrac27@gmail.com', '12345');

INSERT INTO Users(Username, Password, Email, UID)
VALUES('nwchen', 'password', 'nwchen@columbia.edu', '54321');

INSERT INTO Students(Gender, DOB, SchoolYear, CurrentLessonID, UID) 
VALUES('M', '10/27/1997', 2018, '123', '12345');

INSERT INTO Instructors(Gender, DOB, SchoolType, SchoolName, UID) 
VALUES('M', '10/27/1997', 'High School', 'BBCMais', '54321');

INSERT INTO Schools(SchoolName, SchoolID) 
VALUES('BBCMais', '1');

INSERT INTO Classrooms(ClassroomID, InstructorID, Year, Subject, SchoolID)
VALUES('1', '54321', 2018, 'Computer Science', '1');

INSERT INTO Units(UnitName, UnitDescription, UnitID)
VALUES ('Unit', 'Unit', '123');

INSERT INTO Chapters(ChapterID, ChapterName, ChapterDescription, UnitID)
VALUES ('1', 'First Chapter', 'Its the first chapter', '123');

INSERT INTO Lessons(LessonName, LessonContent, LessonID, ChapterID)
VALUES('First Lesson', '{{"hello"}, {"world"}}', '123', '1');

INSERT INTO Lessons(LessonName, LessonContent, LessonID, ChapterID)
VALUES('Second Lesson', '{{"goodbye"}, {"until next time"}}', '1234', '1');

INSERT INTO LessonsCompleted (LessonID, UID)
VALUES ('123', '12345');

INSERT INTO ChaptersCompleted (ChapterID, UID)
VALUES ('1', '12345');