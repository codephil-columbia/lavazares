CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE FUNCTION start() RETURNS void AS $$
    DECLARE
        chapter_0 TEXT := uuid_generate_v4();
        chapter_1 TEXT := uuid_generate_v4();
        chapter_2 TEXT := uuid_generate_v4();
        chapter_3 TEXT := uuid_generate_v4();
        chapter_4 TEXT := uuid_generate_v4();
        chapter_5 TEXT := uuid_generate_v4();
        chapter_6 TEXT := uuid_generate_v4();
    BEGIN
        INSERT INTO chapters (chapterid, chaptername, chapterdescription)
        VALUES 
            -- Chapter Info -- 
            (chapter_0, 'Chapter 0: The Basics', 'In this chapter, you will learn the importance of touch typing as well as proper posture while using the keyboard. You will also encounter the entire finger map of the keyboard, which will serve as your guidelines for the rest of the tutorial.'),
            (chapter_1, 'Chapter 1: Home Row', 'In this chapter, you will learn how to use the home row, which is located on the center of the keyboard. Home row is where you will start your typing journey.'), 
            (chapter_2, 'Chapter 2: Shift and Basic Punctuations', 'In this chapter, you will learn how to use the ''Shift'' keys, which are located on both sides of the keyboard. The ‘Shift’ keys will allow you to access common punctuations. '),
            (chapter_3, 'Chapter 3: Top Row', 'In this chapter, you will learn how to use the top row, which is located above the home row. Mastering the top row allows you to access all English vowels.'), 
            (chapter_4, 'Chapter 4: Bottom Row', 'In this chapter, you will learn how to use the bottom row, which is located below the home row. You have already encountered some of the bottom row keys in Chapter 2, but this chapter will complete your understanding of the bottom row.'), 
            (chapter_5, 'Chapter 5: Number Row', 'In this chapter, you will learn how to use the number row. The number row is important for performing math operations. This chapter will complete your keyboard knowledge.'), 
            (chapter_6, 'Chapter 6: Advanced Content', 'The final chapter of the tutorial will include advanced training exercises. This will bring your WPM and accuracy to the next level, by providing practices through common words, sentences, and paragraphs.');
            

            INSERT INTO lessons (lessonname, lessontext, lessonid, chapterid, minimumscoretopass, lessondescriptions, image)
            VALUES('Lesson 1: The Importance of Touch Typing', 
                '{"Typing faster with better accuracy will help you increase your productivity.", "In this tutorial, you will learn how to touch type. Touch typing is typing without looking at the keyboard to find the keys. If you master touch typing, you will remember the location of keys on the keyboard through muscle memory.", "Touch typing will allow you to type faster with accuracy, increase productivity, and decrease fatigue. Typing can be difficult mentally and physically without touch typing. But learning how to touch type can make typing more enjoyable!"}', 
                uuid_generate_v4(), 
                'e6a18785-98c5-41bc-ad98-ec5d3a243d15', 
                '{0,0,0}', 
                '{"","",""}', 
                '{"","",""}');
            
            INSERT INTO lessons (lessonname, lessontext, lessonid, chapterid, minimumscoretopass, lessondescriptions, image)
            VALUES('Lesson 2: Metrics For Calculating Improvement', 
            '{"How do you know if your typing skills are improving? There are ways to calculate both typing speed and accuracy. In this lesson, you will be introduced to words per minute (WPM) and accuracy.", "Words per minute (WPM) is the number of words you type per minute. The higher your WPM, the faster you are at typing.", "The faster you type, the faster you can communicate with others. You can maximize your time sending an email, researching on the internet, writing project proposals, and more.", "An average person types between 38 and 40 WPM. However, a professional typist can type a lot faster, between 65 and 75 WPM. The record of the fastest English language typist is 216 WPM! Do you think you can beat that?", "Accuracy is the percentage of letters you type that are correct. For example, if you type ‘A’ when you are supposed to type ‘B,’ your accuracy score will decrease.", "Try to keep your average accuracy above 92%! This means you make 8 mistakes for every 100 words typed.", "We will display WPM and accuracy after every exercise you complete. This will allow you to track how well you are doing. You can also refer to “My Progress” page to track your progress over time!"}', 
            uuid_generate_v4(), 
            'e6a18785-98c5-41bc-ad98-ec5d3a243d15', 
            '{0,0,0,0,0,0,0}', 
            '{"", "", "", "", "", "", ""}',
            '{"", "", "", "", "", "", ""}');


            INSERT INTO Users(Username, Password, Email, UID, Occupation, FirstName, LastName) 
            VALUES('tester', '123', 'tester@gmail.com', '12345', 'Student', 'Tester', 'Tester');

            INSERT INTO Students(Gender, DOB, CurrentLessonID, CurrentChapterID, CurrentChapterName, UID)
            VALUES('Male', '10/27/1997', 
            'd3f9c2a3-1edf-42a6-a24d-3a4ad4683036', 'e6a18785-98c5-41bc-ad98-ec5d3a243d15',
            'Chapter 0: The Basics', '12345');

            INSERT INTO Pupils(SchoolYear, UID) VALUES('1st Grade', '12345');
    END;
$$ LANGUAGE plpgsql;

-- DO $$ BEGIN
--     PERFORM start();
-- END $$;

SELECT start();