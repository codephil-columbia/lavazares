CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE FUNCTION start() RETURNS void AS $$
    DECLARE
        chapter_0 TEXT := 'e6a18785-98c5-41bc-ad98-ec5d3a243d15';
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
            
            UPDATE chapters
            SET chapterimage = 'images/chapters/chapter0.svg'
            WHERE chapterid = chapter_0;

            UPDATE chapters
            SET chapterimage = 'images/chapters/chapter1.svg'
            WHERE chapterid = chapter_1;

            UPDATE chapters
            SET chapterimage = 'images/chapters/chapter2.svg'
            WHERE chapterid = chapter_2;

            UPDATE chapters
            SET chapterimage = 'images/chapters/chapter3.svg'
            WHERE chapterid = chapter_3;

            UPDATE chapters
            SET chapterimage = 'images/chapters/chapter4.svg'
            WHERE chapterid = chapter_4;

            UPDATE chapters
            SET chapterimage = 'images/chapters/chapter5.svg'
            WHERE chapterid = chapter_5;

            UPDATE chapters
            SET chapterimage = 'images/chapters/chapter6.svg'
            WHERE chapterid = chapter_6;

            -- Chapter 0 --

            INSERT INTO lessons (lessonname, lessontext, lessonid, chapterid, minimumscoretopass, lessondescriptions, image)
            VALUES
                ('Lesson 1: The Importance of Touch Typing', 
                    '{"Typing faster with better accuracy will help you increase your productivity.", "In this tutorial, you will learn how to touch type. Touch typing is typing without looking at the keyboard to find the keys. If you master touch typing, you will remember the location of keys on the keyboard through muscle memory.", "Touch typing will allow you to type faster with accuracy, increase productivity, and decrease fatigue. Typing can be difficult mentally and physically without touch typing. But learning how to touch type can make typing more enjoyable!"}', 
                    'd3f9c2a3-1edf-42a6-a24d-3a4ad4683036', 
                    chapter_0, 
                    '{0,0,0}', 
                    '{"","",""}', 
                    '{"","",""}')
                ;
            
            INSERT INTO lessons (lessonname, lessontext, lessonid, chapterid, minimumscoretopass, lessondescriptions, image)
            VALUES
                ('Lesson 2: Metrics For Calculating Improvement', 
                    '{"How do you know if your typing skills are improving? There are ways to calculate both typing speed and accuracy. In this lesson, you will be introduced to words per minute (WPM) and accuracy.", "Words per minute (WPM) is the number of words you type per minute. The higher your WPM, the faster you are at typing.", "The faster you type, the faster you can communicate with others. You can maximize your time sending an email, researching on the internet, writing project proposals, and more.", "An average person types between 38 and 40 WPM. However, a professional typist can type a lot faster, between 65 and 75 WPM. The record of the fastest English language typist is 216 WPM! Do you think you can beat that?", "Accuracy is the percentage of letters you type that are correct. For example, if you type ‘A’ when you are supposed to type ‘B,’ your accuracy score will decrease.", "Try to keep your average accuracy above 92%! This means you make 8 mistakes for every 100 words typed.", "We will display WPM and accuracy after every exercise you complete. This will allow you to track how well you are doing. You can also refer to “My Progress” page to track your progress over time!"}', 
                    uuid_generate_v4(), 
                    chapter_0, 
                    '{0,0,0,0,0,0,0}', 
                    '{"", "", "", "", "", "", ""}',
                    '{"", "", "", "", "", "", ""}')
                ;

            INSERT INTO lessons (lessonname, lessontext, lessonid, chapterid, minimumscoretopass, lessondescriptions, image)
            VALUES 
                ('Lesson 3: Body Posture', 
                    '{"Sit up as straight as possible and keep your back straight.", "Rest your arms on the edge of the table. You should be able to freely move your wrists and hands.", "Keep your arms, shoulders, neck, and back relaxed."}', 
                    uuid_generate_v4(), 
                    chapter_0, 
                    '{0,0,0}', 
                    '{"","",""}', 
                    '{"images/lessons/proper_body_position.svg","images/lessons/proper_body_position.svg","images/lessons/proper_body_position.svg"}')
                ;


            INSERT INTO lessons (lessonname, lessontext, lessonid, chapterid, minimumscoretopass, lessondescriptions, image)
            VALUES
                ('Lesson 4: Hand Posture', 
                    '{"Make sure your hands are raised so that your palms are not touching the keyboard or table.", "Make sure your fingers are curved and pointed down at the keys."}', 
                    uuid_generate_v4(), 
                    chapter_0, 
                    '{0,0}', 
                    '{"",""}', 
                    '{"images/lessons/proper_hand_position.svg","images/lessons/proper_hand_position.svg"}')
                ;

            INSERT INTO lessons (lessonname, lessontext, lessonid, chapterid, minimumscoretopass, lessondescriptions, image)
            VALUES
                ('Lesson 5: Keyboard and Finger Maps', 
                    '{"In this tutorial, we will refer to your fingers as pinky, ring, middle, pointer, and thumb. Please see the chart below.","Please see the color-coded keyboard below. This keyboard will help you visualize which finger is used on each key.", "Always return your fingers to the starting position (`ASDF` for the left hand and `JKL&#59` for the right hand) as shown below.", "Always imagine this keyboard layout while keeping your eyes at the screen.", "Use the thumb of your dominant hand to press the `Spacebar`.", "Keep practicing with this keyboard layout. Even though it may seem difficult at first, you will be able to type easily and quickly after the tutorial is over."}', 
                    uuid_generate_v4(), 
                    chapter_0, 
                    '{0,0,0,0,0,0}', 
                    '{"","", "", "", "", ""}', 
                    '{"images/lessons/hand_label_map.svg","images/lessons/keyboard_layout_map.svg","images/lessons/keyboard_layout_map.svg","images/lessons/keyboard_layout_map.svg","images/lessons/keyboard_layout_map.svg","images/lessons/keyboard_layout_map.svg"}')
                ;

            INSERT INTO lessons (lessonname, lessontext, lessonid, chapterid, minimumscoretopass, lessondescriptions, image)
            VALUES
                ('Lesson 6: Ways to Improve your Typing Skills', 
                    '{"Practice, practice, practice! Practicing is the most effective way to improve your performance in anything, including typing! Even just practicing 30-minutes a day can help you improve your typing speed and accuracy.","Type without looking at the keyboard! Even when you are practicing alone at home, motivate yourself to type without looking at the keyboard. This will allow you to develop a habit of typing faster.","Slow down! In the beginning, try to focus on accuracy. The speed will follow when you practice using this tutorial. To increase accuracy, try not to use the ‘Delete’ key. Using the ‘Delete’ key teaches you how to fix errors, not how to reduce errors.","Find a rhythm! You should establish and maintain a rhythm while typing. This means that keystrokes should come at equal intervals. Rhythm is important because it helps you improve your accuracy."}', 
                    uuid_generate_v4(), 
                    chapter_0, 
                    '{0,0,0,0}', 
                    '{"","","",""}', 
                    '{"","","",""}')
                ;


            -- Chapter 1 --
            INSERT INTO lessons (lessonname, lessontext, lessonid, chapterid, minimumscoretopass, lessondescriptions, image)
            VALUES
                ('Lesson 1: Introduction', 
                    '{"Welcome to your very first chapter! In Chapter 1, you will learn one of the most important concepts in typing: home row.", "Home row is the middle horizontal row of the keyboard. It is where your fingers return to when you are not typing. Your fingers should lightly keep in touch with their “homes” in the home row. This will help you find a reference point for other keys and allow you to type without looking at your hands.", "Home row includes keys: ‘A, S, D, F, J, K, L, semicolon (&#59)’.", "The image below illustrates where to place your fingers on the home row. Place your fingers gently on their respective keys, but make sure your fingers are lifted so that you are not actually pressing them!", "You can always find home by placing your fingers on the small bumps on ‘F’ and ‘J’."}', 
                    uuid_generate_v4(), 
                    chapter_1, 
                    '{0,0,0,0,0}', 
                    '{"","","","",""}', 
                    '{"images/lessons/home_row_map.svg","images/lessons/home_row_map.svg","images/lessons/home_row_map.svg","images/lessons/home_row_map.svg","images/lessons/home_row_map.svg"}')
                ;

            INSERT INTO lessons (lessonname, lessontext, lessonid, chapterid, minimumscoretopass, lessondescriptions, image)
            VALUES
                ('Lesson 2: Left Hand (ASDF)', 
                    '{"The home row keys for the left hand are ‘A, S, D, F’.", "Type the ‘F’ key using your left pointer.", "For spaces, use whichever thumb you feel most comfortable with. Remember, your thumbs should always remain at the ‘Spacebar’!", "Type using the ‘F’ and ‘Spacebar’ below.", "Great! Now let’s add the rest of the left hand. Please see the chart below to learn which fingers are used for which keys. Use your middle for ‘D’, ring for ‘S’, and pinky for ‘A’.", "Let’s practice some more with the ‘D’ key!", "Great! Now let’s add the ‘S’ key!", "Awesome job! Now let’s add the final key: ‘A’.", "Great job so far! As you complete the next exercise, try not to look down at the keyboard when you type.", "Awesome job! Let’s do another exercise!", "Great job so far! Let’s try another exercise.", "One last exercise! Let’s go!"}', 
                    uuid_generate_v4(), 
                    chapter_1, 
                    '{0,60,0,60,0,60,60,60,0,60,60,60}', 
                    '{"","ffffffffff", "", "f f ff fff ff f f f f f f ff fff f f ff ff fff ff fff", "", "d d dd fd df f dd dddf fd df fd dd df fd dff fd ddf dddf d df dfd fdf d fddd dfd", "ssd sfd sdf dfs ssf fds sfd dss sfd fsd sd ssd s sfds s dfs s s s d f df s f d s", "add das sdas aad fas ads afa fad das fads aad asda asa a d aa d dd a aa", "", "ffff aaaa dddd ssss fa sd adad fs da af fs dsdf fdsa fdsa fads fads fff asa as a s d f f d s a aa", "f f ffd ddfs a aaf s a d asdf fdsa a s d f f d s a asdf d af d af ddf ffd aaf ssd dda ffa dfd ada sfs fsf dad dads fads", "ffa dda ada dad ada fsd fsd daf fsd fas sfd ads add dad dda ada ffs ssf dff ffd dad dds fsd ads adf ads adf"}', 
                    '{"","","","","","","","","","","",""}')
                ;

            INSERT INTO lessons (lessonname, lessontext, lessonid, chapterid, minimumscoretopass, lessondescriptions, image)
            VALUES
                ('Lesson 3: Right Hand (JKL&#59)', 
                    '{"Awesome work on that left hand! Now let’s find a home for the your right hand fingers.", "Once again, make sure you gently place your right pointer on the ‘J’ key. The home row keys for the right hand are ‘J, K, L, &#59 (semi-colon)’.", "Type the ‘J’ key using your right pointer.", "Great! Now let’s add the rest of the home row for the right hand. Please see the chart below to learn which fingers are used to type which keys. Use your middle for ‘K’, ring for ‘L’, and pinky for ‘&#59’.", "Let’s practice some more with the ‘K’ key!", "Great! Now let’s add the ‘L’ key!", "Awesome job! Now let’s add the final key ‘&#59’.", "You are doing great! Let’s try an exercise with your entire right hand.", "Great job! Let’s do one more exercise."}', 
                    uuid_generate_v4(), 
                    chapter_1, 
                    '{0,0,60,0,60,60,60,60,60}', 
                    '{"","", "jjj j j jj jj jjj j", "", "jk kj kkj kj jk jkk kk k kk kk jk kjk kkjkj jjk kj jkkj k", "lkj jlkj llj lkl jlk ljlkj l l jlkjl jlkj jlkj llk jlkl lklj jlkjl l llkl ljkl lj l", "&#59 &#59&#59 k&#59l jlj&#59 &#59kl &#59j&#59 &#59ll&#59k j&#59kl&#59 &#59k&#59 &#59&#59 &#59k&#59j&#59lk j&#59l&#59 k&#59 k&#59 &#59kk &#59&#59kl&#59 j&#59 lk&#59 j&#59 lk&#59j &#59", "j j jjk k k kkl l l ll&#59 &#59 &#59 &#59&#59j j&#59 kl lk j&#59 l k jjj kkk jjj kkk lll kkk lll jjj &#59&#59&#59 jjl llj jjl kkj jkj kjk jkj jkj", "xlll &#59&#59&#59 kkk j&#59l jlk k&#59l jlk jlk k&#59l lj&#59 &#59kl &#59jl &#59&#59&#59 j k l &#59 &#59 l k j j&#59 &#59j kl lk k&#59 &#59k l&#59 &#59l jk&#59"}', 
                    '{"","", "","","","","","",""}')
                ;

            INSERT INTO lessons (lessonname, lessontext, lessonid, chapterid, minimumscoretopass, lessondescriptions, image)
            VALUES
                ('Lesson 4: Left and Right Hand (ASDFJKL&#59)', 
                    '{"Great job so far! Now let’s have some practice with your left and right hand together.", "Don’t forget all of the proper typing techniques that you have learned so far.", "Type the letters below.", "Type the letters below."}', 
                    uuid_generate_v4(), 
                    chapter_1, 
                    '{0,60,60,60}', 
                    '{"","fff jjj jjj ddd ddd kkk kkk sss sss lll lll aaa aaa &#59&#59&#59 aad jjl llk ddk &#59&#59s aak ssl lld ffj jjf alsk fjdk a&#59sl sldk","a s d f j k l &#59 &#59 a l s k d j f fdf fdf fdf dfd sas sas ssa jkj lkl jlj llk asdfjkl&#59 fjdk fjfj lsls add daa add dda llj jjl jjl llj","ds dasf&#59 llfaaaa jlakaflaf kja ajlldaskdff&#59 al kkllkjfk afaldjskdsdfafakfld fjslfsljsllad sdsjsallklk&#59 lfdjddjkd&#59 ajfsslddjs dlksjjjafj&#59 ajkadklf"}', 
                    '{"","","",""}')
                ;

            INSERT INTO lessons (lessonname, lessontext, lessonid, chapterid, minimumscoretopass, lessondescriptions, image)
            VALUES
                ('Lesson 5: Extended Home Row (ASDFGHJKL&#59)', 
                    '{"You might have noticed that we have not learned two letters on the home row yet. Can you identify them?", "We are missing ‘G’ and ‘H’! The keys we covered so far can all be typed without moving your fingers. Now, let’s learn how to move your fingers to reach ‘G’ and ‘H’.", "To reach ‘G’ and ‘H’, you will need to move your pointers. Try practicing this motion by following the illustration below. Your left pointer will hit the ‘G’ key and the right pointer will hit the ‘H’ key.", "Don’t forget to move your fingers back to the home row position after you type ‘G’ and ‘H’ keys.", "Type the keys below.", "Type the keys below."}', 
                    uuid_generate_v4(), 
                    chapter_1, 
                    '{0,0,0,60,60,60}', 
                    '{"","","","gg gg hh hh gh hg gh hg h h g g h h g fg fg gf fg gfdsa asdfg gf fg dfg gfd g f g f g a f s g f g", "hj hj jh hj hjkl&#59 &#59lkjh hj hk hl lh hj hjk &#59 &#59kjh hjkl&#59 h j h j k gh hg jf fj gj hf gk hd dj fl a&#59 lh hj ga sg sh lg dh gal lad", "hhh hhh ggg ggg hhh ggg jjj fff jjj ggg hhh fj jfj fjf hgh ghg hfh gjg jgj j f gjg hf gj hg gh hf jg fh ghg jfj"}', 
                    '{"","","images/lessons/home_row_map.svg","","",""}')
                ;

            INSERT INTO lessons (lessonname, lessontext, lessonid, chapterid, minimumscoretopass, lessondescriptions, image)
            VALUES
                ('Chapter 1 Test', 
                    '{"Try to type these letters as accurately as possible!","Remember to keep your hands on the home row!","Type the following words and letters. Try not to look down at the keyboard and remember your home!"}', 
                    uuid_generate_v4(), 
                    chapter_1, 
                    '{70,70,70}', 
                    '{"asdfghjkl&#59 &#59lkjhgfdsa alfalfas&#59 ashfalls haggadah&#59 haggadas halakahs&#59 halakhas halalahs", "hask khafk ah lads lag kaf fash&#59 flag flaks&#59 kafs alfs as lash sad&#59 ash ask flags&#59 flak dahs jag jags&#59 lags lakhs la khafs&#59 lad jaks dahls&#59", "flags flash shall gaff jag&#59 ha flag half sagas&#59 gags gash salad jags glad&#59 lash hash has as dad fad haha&#59 haj asks alls lass lash&#59 shag&#59 salad slags glass flags daffs algas flash add lad"}', 
                    '{"","",""}')
                ;


            -- Chapter 2 --
            INSERT INTO lessons (lessonname, lessontext, lessonid, chapterid, minimumscoretopass, lessondescriptions, image)
            VALUES
                ('Lesson 1: Introduction', 
                    '{"Great job so far! You are well on your way to becoming an excellent typist!","But before we move onto the other rows, you need to understand one of the most important keys in typing. The ‘Shift’!"}', 
                    uuid_generate_v4(), 
                    chapter_2, 
                    '{0,0}', 
                    '{"",""}', 
                    '{"",""}')
                ;

            INSERT INTO lessons (lessonname, lessontext, lessonid, chapterid, minimumscoretopass, lessondescriptions, image)
            VALUES
                ('Lesson 2: Introduction to Shift Key', 
                    '{"Most keyboards have two ‘Shift’ keys: one on the left and one on the right.","The ‘Shift’ key is pressed using the pinky. Please see the illustration below and press the ‘Shift’ key to continue.", "Don’t forget to return your pinky back to Home Row after you release the ‘Shift’ key.", "‘Shift’ keys can be tricky at first. Try to use the ‘Shift’ key that is on the opposite side of the key you are typing.", "For instance, press the left ‘Shift’ key when you are typing a key on the right side of the keyboard.", "And press the right ‘Shift’ key when you are typing a key on the left side of the keyboard."}', 
                    uuid_generate_v4(), 
                    chapter_2, 
                    '{0,0,0,0,0,0}', 
                    '{"","","","","",""}', 
                    '{"images/lessons/bottom_row_map.svg","images/lessons/bottom_row_map.svg","images/lessons/bottom_row_map.svg","images/lessons/bottom_row_map.svg","images/lessons/bottom_row_map.svg","images/lessons/bottom_row_map.svg"}')
                ;

            INSERT INTO lessons (lessonname, lessontext, lessonid, chapterid, minimumscoretopass, lessondescriptions, image)
            VALUES
                ('Lesson 3: Using Shift with Home Row Letters (A, S, D, F, H, J, K, L, &#59)', 
                    '{"Now it is time to use the ‘Shift’ key with all of the home row keys!","The ‘Shift’ key serves two functions. First, it makes a lowercase letter an uppercase letter. So if you press on both ‘Shift’ and ‘f’ at the same time, it gives you ‘F’. For non-letter keys, it types the character that is located on top. For instance ‘&#59’ used with a ‘Shift’ key, will turn into ‘:’. Please see the illustration below.","Let’s try a few more exercises using the ‘Shift’ key!","Great job so far. Now let’s practice on actual words.","Awesome job! Let’s do another practice!"}', 
                    uuid_generate_v4(), 
                    chapter_2, 
                    '{0,0,70,70,70}', 
                    '{"","","AaA SsS DdD FfF GgG hHh jJj kKk lLl &#59:&#59 aSdFgHjKl: &#59 L K J H G F D S A aaa AAA aaa AAA &#59&#59&#59 ::: &#59&#59&#59 ::: aA aS Aa Ss :&#59 &#59: Ll lL LLJ FfA AaF fFA Faf AFA llJ JJL","Dallas ADA flag Alaska Flask AS : jads Hash sad Lass : Dads : Haha :salad&#59 salsa: Halls as&#59 gas: Add Flasks, ash: Alaska Lass &#59 saga: half flag:","dash Fads lads Lags dals Dahs lash Flag Shad dhal gals Fags Hags gads dags Gash Flak"}', 
                    '{"","images/lessons/how_shift_works.svg","","",""}')
                ;

            INSERT INTO lessons (lessonname, lessontext, lessonid, chapterid, minimumscoretopass, lessondescriptions, image)
            VALUES
                ('Lesson 4: Shift with &#59, :, ‘, “', 
                    '{"For our next lesson, we will explore how to use the ‘Shift’ key with the semicolon (&#59) to produce a colon (:) and how to use it with an apostrophe (’) to produce a quotation (“).","To reach these special characters, you will need to extend your pinky. Just like you did to reach the keys ‘G’ and ‘H’, you need to extend your right pinky to reach the the apostrophe (‘).","Let’s try reaching for the apostrophe (‘) a few times.","Great! Now let’s add the ‘Shift’ key to convert the semi-colon (&#59) to a colon (:) and the apostrophe (‘) to a quotation mark (“)."}', 
                    uuid_generate_v4(), 
                    chapter_2, 
                    '{0,0,70,70}', 
                    '{"","","‘ &#59’ ‘&#59’ &#59’ l&#59’ ‘&#59k ‘’&#59’ ‘ “add” “jlka”","‘ “ &#59&#59: “ ‘ “ ‘ “ “ “ : : : &#59 ‘&#59 l k” &#59:’l “ &#59 “ k’:”l&#59"}', 
                    '{"images/lessons/special_key_explanation.svg","","",""}')
                ;

            INSERT INTO lessons (lessonname, lessontext, lessonid, chapterid, minimumscoretopass, lessondescriptions, image)
            VALUES
                ('Lesson 5: Basic Punctuations (,.!?)', 
                    '{"We have now learned everything on the home row! Congratulations!","Before we move on to other letters, let’s cover how to properly reach for other basic punctuations.","If you look at the illustration of the keyboard below, you can see that the comma (,), period (.), and question mark (?) are located on the bottom row. You can access these keys using your middle, ring, and pinky.","You may have noticed that you need to use the ‘Shift’ key to access the question mark (?). To do this, hold the ‘Shift’ key and press the (/) key.","Now for the exclamation point (!), your left pinky will need to travel two rows up from the home row to press the ‘1’ key. But make sure to hold the ‘Shift’ key with your right pinky before you press the ‘1’ key.","Great! Now let’s put everything together. Don’t forget to return your fingers to your home row after you move them away to press other keys!","Great job! Let’s try one more exercise for this lesson."}', 
                    uuid_generate_v4(), 
                    chapter_2, 
                    '{0,0,0,70,70,70,70}', 
                    '{"","","","/? ???/ ?/ ?/? /?/ /&#59? ?&#59/’/”","a!aaA A!A!a1 a!1!a!a A!!","Salad! Dad’s “haha”? Alaska, flag. Add, salsa’s: lass HAHA?! “Dallas” Alfalfas, ashfalls&#59 “haggadah,” Haggadas! halakahs, halakhas? Halalahs.","flags dhaks Flash! flask? Lakhs. flaks “Skald” dhals. glads, dahls, slag? dahs: flak&#59 Gash Dhak hadj Dash! Daks lags, dhal: lakh “hags” fags, fash! Half"}', 
                    '{"","images/lessons/keyboard_layout_map.svg","images/lessons/keyboard_layout_map.svg","","","",""}')
                ;

            INSERT INTO lessons (lessonname, lessontext, lessonid, chapterid, minimumscoretopass, lessondescriptions, image)
            VALUES
                ('Lesson 6: Caps Lock', 
                    '{"What if you want to type multiple uppercase letters in a row? While you can press the ‘Shift’ key each time, you can also use the ‘Caps Lock’ key!","When pressed, the ‘Caps Lock’ key allows all letters to be generated in capitals until deactivated.","Usually, ‘Caps Lock’ only works on letter keys. For instance, if you press a non-letter key such as ‘1’ with your ‘Caps Lock’ on, you will still get ‘1’ instead of ‘!’.","Don’t forget to turn off ‘Caps Lock’ when you are done using it."}', 
                    uuid_generate_v4(), 
                    chapter_2, 
                    '{0,0,0,70}', 
                    '{"","","","ASDFLKJGH LKJ ALSDF. ASLKJG. ASLGK’S “LJAKJSFLD”! ASLG KSJLKJFA. SAD LAG! DAH GAD ‘HAD’ aSH “LAH” SHA! FAH GAS, GAL, dAG. DAK: KAF DAL LAD FAD hAJ ASK SAL, SAG? JAG FAG HAG JAk SKa"}', 
                    '{"","","",""}')
                ;

            INSERT INTO lessons (lessonname, lessontext, lessonid, chapterid, minimumscoretopass, lessondescriptions, image)
            VALUES
                ('Lesson 7: Delete Key', 
                    '{"Hopefully, you didn’t need to use the ‘Delete’ key that much so far! However, it is still important to review how to use the ‘Delete’ key.","Use your right pinky to reach for the ‘Delete’ key. But don’t forget to return your right pinky back home after you are done deleting!","Instead of relying on the ‘Delete’ key when you type, try to type carefully and accurately so that you do not have to use the ‘Delete’ key. Making mistakes is costly because you have to delete it then retype it. Just get it right the first time!"}', 
                    uuid_generate_v4(), 
                    chapter_2, 
                    '{0,0,0}', 
                    '{"","",""}', 
                    '{"","",""}')
                ;

            INSERT INTO lessons (lessonname, lessontext, lessonid, chapterid, minimumscoretopass, lessondescriptions, image)
            VALUES
                ('Chapter 2 Test', 
                    '{"Focus on accuracy during this test. Make sure to get all of the characters correct!","Great job! Let’s do one more exercise to complete this chapter!"}', 
                    uuid_generate_v4(), 
                    chapter_2, 
                    '{80,80}', 
                    '{"“ha FLAG half lash&#59” gags ‘gash?’ salad: jags! “glad&#59 LaSH hash” has as dad, fad? haha&#59 “haj ASKS!” alls lass? Shag&#59 flags. “Flash, Shall gaff!” jagsalad? slags/ GLASS! FLAGS DAFFS AlGas “flash” add lad sagas&#59","S&#59 lldk! Kdgd fKhd ?gA k?&#59, “!” kFl “sl” ??dhh&#59f&#59 Kgk.l jJaa Fdjg? DgLgG ,dlh&#59, “,”dfGh&#59J ,gd a?sJl j!gd aSjh dj “!g &#59,kLdd&#59 gf,ka flaks skald dhals glads dahls"}', 
                    '{"",""}')
                ;


            -- Chapter 3 --
            INSERT INTO lessons (lessonname, lessontext, lessonid, chapterid, minimumscoretopass, lessondescriptions, image)
            VALUES
                ('Lesson 1: Introduction', 
                    '{"Congratulations! You have now graduated from the home row! We will now move onto the top row.","Top row refers to the row above your home row. It includes the keys ‘Q, W, E, R, T, Y, U, I, O, P’.","Don’t forget to return your fingers back to your home row after you type something on the top row."}', 
                    uuid_generate_v4(), 
                    chapter_3, 
                    '{0,0,0}', 
                    '{"","",""}', 
                    '{"","images/lessons/top_row_map.svg",""}')
                ;

            INSERT INTO lessons (lessonname, lessontext, lessonid, chapterid, minimumscoretopass, lessondescriptions, image)
            VALUES
                ('Lesson 2: Left Hand (QWERT)', 
                    '{"Keys ‘Q, W, E, R, T’ on the top row can be reached with your left hand.","You will be moving your fingers much more during this chapter. Remember to keep your fingers curved and pointed down on the keys.","Please see the chart below to learn which fingers are used to type which keys. Use pinky for ‘Q’, ring for ‘W’, middle for ‘E’, pointer for ‘R’, and pointer again for ‘T’!","Let’s first add the letter ‘Q’! Make sure to slow down to reach the new keys with accuracy.","ow let’s add the letter ‘W’!","Time to practice with ‘E’!","Don’t forget about ‘R’!","Finally, time to practice with ‘T’.","Great job! That was a lot of new keys to learn. Before we move on, let’s review some good typing habits. Remember to have the proper body posture when you type. In the next screen, we will review some of these tips for proper body posture.","Sit up as straight as possible and keep your back straight.","Rest your arms on the edge of the table. You should be able to freely move your wrists and hands.","Keep your arms, shoulders, neck, and back relaxed.","Let’s practice with all of the new keys that you have learned!","Awesome job! Let’s do one more practice!"}', 
                    uuid_generate_v4(), 
                    chapter_3, 
                    '{0,0,0,80,80,80,80,80,0,0,0,0,80,80}', 
                    '{"","","","Qqjda qfjq? Qhgdfgd! Qqq qhglaj! qgqf Qkga? Q, Qfga, qhh q qkj Qflfqhk “Qgqk” qhfjlg qga.","Wfgl, Wk wjjd Wwlwh Wdgg? Wfqagwh “Wwsfsja” wfq, Whdslh? Wklwjl, wqghg. Wjk wfj Wjas “Wqhddq” ws.","Eqhwq Eejgj eajhkas “Ef” ejqhhka, ejfga, efeqlq! “Ehqekaj” Ees: Egh Ejaha! ekeawjf efeed","Rswsqgd raq rw Rdkqsw, Rkd, rwrwef rgfhR rwwqd. Rkawagr! Rrdfj? Rses! Rjaerg? Rqaj! Rhrhdfr ‘rqgdekl’ rrdqfa: RrRj","Tjqagfh twge twsrgs tkdhqr, Taswh. Twwefs. tqdhekf! Thlgksl “tglr” tewj Twtrdrr, Tllfer, Tfwrgeg, tfhea Tdewafa","","","","","QAfd! Raws rare fares? Wrastled warstled Dwarfest selfward! “Leftward” wreath jehads “Father” staler THAWED larked staled fadges, drakes, grafts, talked.","awa, awash awl, awls ‘daw,’ dawk,ale, Alee! alef. alefs, ales “algae,” Grad grads Grass Haar haars, haggard, haggards, harass, hard, hards, hark, harks harl Harls harsh, jagra Jagras jar. jarl, jarls? jarrah, jarrahs, hafts, halt, halts. Hast!"}', 
                    '{"images/lessons/top_row_map.svg","images/lessons/top_row_map.svg","images/lessons/top_row_map.svg","","","","","","","images/lessons/proper_body_position.svg","images/lessons/proper_body_position.svg","images/lessons/proper_body_position.svg","",""}')
                ;

            INSERT INTO lessons (lessonname, lessontext, lessonid, chapterid, minimumscoretopass, lessondescriptions, image)
            VALUES
                ('Lesson 3: Right Hand (YUIOP)', 
                    '{"Now time to learn the right-hand side of the top row: ‘Y, U, I, O, P’!","You will be moving your fingers much more during this chapter. Remember to return your fingers back to the home row when you are not using them.","Please see the chart below to learn which fingers are used to type which keys. Use the pinky for ‘P’, ring for ‘O’, middle for ‘I’, pointer for ‘U’, and pointer again for ‘Y’.","Let’s first add the letter ‘P’! Don’t forget to slow down to accurately reach the new key.","Great work! Let’s add the key ‘O’.","Time to learn the key ‘I’.","Let’s add ‘U’!","The final key ‘Y’. Let’s practice it!","Great! You have now learned all of the top row. Let’s practice some more!","One last exercise to finish the lesson!"}', 
                    uuid_generate_v4(), 
                    chapter_3, 
                    '{0,0,0,80,80,80,80,80,80,80}', 
                    '{"","","","Padshahs pah pal Pall palls palp palpal Palps Pals pap papa papal papas Paps pas pash Pasha pashas","oaf oafs Oak oaks Odd Odds ods of off offal Offals offload offloads Offs ogdoad ogdoads oh oho Ohs Oka okas Old Olds olla ollas ooh","id Ids if iff Ifs Ilia iliad iliads ilial Ilk Ilka ilks ill ills Is Ia Ita Ipoew Iqrt Ipq Iwpee","Ugalis Uplaid Uhs up upas upload ugh Uphold upholds Ughs Upgo upfolds uphild ugs uploads Updo updos uh upfold","Yogi Yuko You youk youks yah Yahs Yak yous yuga yolks Yugas yugs Yuk yukos Yuks yup yogis yoks yokul yold yolk","Oaf Yap padouk Yak yaks uploads, pads, Paid, paik oafs Paiks “oaks” Ups upsy? Oaky odal Ya yads Oak us yald.","Shipload “goldfish” ladyfish, FLAGSHIP, oafishly, Ladyship. aguishly! haploidy? “haploids” Holidays!! Hidalogs uphods soapily dishful huskily"}', 
                    '{"","","images/lessons/top_row_map.svg","","","","","","",""}')
                ;

            INSERT INTO lessons (lessonname, lessontext, lessonid, chapterid, minimumscoretopass, lessondescriptions, image)
            VALUES
                ('Lesson 4: Left and Right Hand (QWERTYUIOP)', 
                    '{"We will now combine your left and right hands for the top row. Remember to return your fingers back to the home row when you are not using them.","Make sure to slow down to reach the new keys with accuracy.","Great job! That was a lot of new keys to learn. Before we move on, let’s review some good typing habits. Remember to have the proper hand posture when you type! In the next screen, we will review some of these tips for proper hand posture.","Make sure your hands are raised so that your palms are not touching the keyboard or table.","Make sure your fingers are curved and pointed down at the keys.","Let’s practice some more!","Type the words below."}', 
                    uuid_generate_v4(), 
                    chapter_3, 
                    '{0,80,0,0,0,80,80}', 
                    '{"","Quatrefoils Outsparkled! “playwrights” Fluoridates, profligates, righteously. Waterdogs drawliest ‘outprayed’ forestial: laughters: horsetail? euphorias odalisque UPLIFTERS playhouse, Upgrowth, sideograph.","",",","","Dihwater headfirst waterlogs PREADUITS uplighted! autopsied redisplay saprolite plowhead skryolites foresight, preflight, holidayer? afterglow. horseplay ‘doughlike’ “polyhedra,” Wordplays! Thioureas? Leadworks ROPEWALKS","Pyet typier pyot Query Ewt eyot ire it outer Outre Owe ower pyrite Qi quep typer quey Pyre Typo Ow tyre tyro Opt opter Uey up eutropy Upter uptie"}', 
                    '{"","","","images/lessons/proper_hand_position.svg","images/lessons/proper_hand_position.svg","",""}')
                ;

            INSERT INTO lessons (lessonname, lessontext, lessonid, chapterid, minimumscoretopass, lessondescriptions, image)
            VALUES
                ('Chapter 3 Test', 
                    '{"Focus on accuracy during this test. Make sure to get all of the characters correct!","Complete the final typing exercise for this chapter!"}', 
                    uuid_generate_v4(), 
                    chapter_3, 
                    '{90,90}', 
                    '{"Gatefolds, dialogers! “Atrophies,” giftwares? Outfields, Walkyries, Headworks, Frugality! Dishtowel “daughters” ‘grapeshot: fieldwork deflators ‘sprightly’ prudishly? Polarised! Therapsid? “Jailhouse” Southerly. Filatures! “Hairstyle” hysteroid Doughtily outglares outglared! Sailer sailed? “Lugers” widget shriek? Guyots. Ligers Golder.","Epopee, equity! euripi? Irrupt ‘orrery,’ Output truer tuque tutee, tutor, tutti, tutty tuyer Tweet, twerp, “twier,” twirp pewterer, portiere potterer? eyepopper, Outpourer pirouette, potpourri preterite Preppier, preterit prettier priority: properer!"}',
                    '{"",""}')
                ;


            -- Chapter 4 --
            INSERT INTO lessons (lessonname, lessontext, lessonid, chapterid, minimumscoretopass, lessondescriptions, image)
            VALUES
                ('Lesson 1: Introduction', 
                    '{"Great job so far! Now that you have mastered the home and top row, there is only one more row to master. Time to learn the bottom row!","Bottom row refers to the row below your home row. It includes the keys ‘Z, X, C, V, B, N, M’.","Just like you did for the top row, don’t forget to return your fingers back to your home row after you type something on the bottom row."}', 
                    uuid_generate_v4(), 
                    chapter_4, 
                    '{0,0,0}', 
                    '{"","",""}', 
                    '{"","images/lessons/bottom_row_map.svg",""}')
                ;

            INSERT INTO lessons (lessonname, lessontext, lessonid, chapterid, minimumscoretopass, lessondescriptions, image)
            VALUES
                ('Lesson 2: Left Hand (ZXCVB)', 
                    '{"Keys ‘Z, X, C, V, B’ on the bottom row can be reached with your left hand.","Please see the chart below to learn which fingers are used to type which keys. Use pinky for ‘Z’, ring for ‘X’, middle for ‘C’, pointer for ‘V’, and pointer again for ‘B’!","This exercise will introduce the ‘Z’ key. Let’s start!","Great! Now let’s add ‘X’.","Excellent job so far. Let’s add key ‘C’ now.","Almost finished with the bottom row! Time to add ‘V’.","Last key! Let’s practice the ‘B’ key.","Awesome! We finished the bottom row for the left hand. Before we complete this lesson, let’s remember review some tips to keep in mind when typing to improve your skills.","Find a rhythm! You should establish and maintain a rhythm while typing. This means that keystrokes should come at equal intervals. Rhythm is important because it helps you increase accuracy.","Let’s practice with all of the new keys that you have learned!","Great job! Remember to find a good rhythm when typing this time.","Awesome! One last exercise to complete this lesson."}', 
                    uuid_generate_v4(), 
                    chapter_4, 
                    '{0,0,90,90,90,90,90,0,0,90,90,90}', 
                    '{"","","Zags zero “Zoo” zealot, zo zoril. Zoea!l Zaire zite Zati zeta&#59 zila “zoea” zit ‘zoa’ Zerda? zed&#59 zel! zea? Zelator zeal","Xi xu “Xrayed” xystoi xrays! xyloid, xylose. Xray xyst? Xyster! xysti’ ‘Xis’","Capsulized! captious. captiously “Captured” cardy captures capul capuls “Caput” car Card, Cardhouse, captor, captors. cardi? cardie “cardies”","Volts voluspa! Vorpal Vortex vortical, vortices volute voluted volutes Votaries, Votary, vote Voted, voter! “Voracity” Vorlage, vorlages.","Backstory Bagels? backveld bailouts bails. bait Bade badge bagful Backswept backsword Bafts, bag, bagel. backup! Backups badger.","","","Cob Voice cozie, Beau, Cub cube! beaux Bivouac obia Ova ox cobza Cove. zoea bize Boa Box coxae “Coze” vie Zoa&#59 zax zobu buaze! cue cuz vocab voe zebu.","Cauterize Ectozoa. Ecotour! ebriate&#59 bioactive Bacterize. Activizer Exoteric exorcize Obviator, Icewater, exuviate? Voicebox. trivia Towier towbar Vizirate equator victoria.","Heartblock Hawseblock cratefuls? Crashdove crashdive. Hardcopies grubstaked Capsulized, Caprifoles grouchiest harpylike “Butlership” bushwalker Harpsicle. Halicores? calfdozers Cabriolets byproducts."}', 
                    '{"","images/lessons/bottom_row_map.svg","","","","","","","","","",""}')
                ;

            INSERT INTO lessons (lessonname, lessontext, lessonid, chapterid, minimumscoretopass, lessondescriptions, image)
            VALUES
                ('Lesson 3: Right Hand (NM,./)', 
                    '{"You are almost done learning all of the letters on the keyboard! Please see the chart below to learn which fingers are used to type which keys. Use the pinky for forward-slash (/), ring for period (.), middle for comma (,), pointer for ‘M’, and pointer again for ‘N’.","Since you are already familiar with the period (.) and comma (,), we will just focus on forward-slash (/), ‘M’, and ‘N’.","Complete the exercise below for the ‘N’ key.","Great! Now let’s practice with the ‘M’ key.","Final key! Forward slash (/) is not used often, but it is still useful.","Great job! Since you have already learned the period (.) and comma (,), we are all done! But before we continue, let’s review good typing habits to improve your skills.","Set certain goals based on your current typing speed and try to aim faster each session you practice. By motivating yourself to hit certain targets, you will witness your typing speed improve significantly over time.","Let’s practice with all of the new keys that you have learned with your right hand."}', 
                    uuid_generate_v4(), 
                    chapter_4, 
                    '{0,0,90,90,90,0,0,90}', 
                    '{"","","Nag nags? neighs noils Naif “nod” nodal negs Naifs Noise nose? naiks Neif neifs Neig, Nodalise node! nail No.","“Medial” Mensh Median meloids, Mid. midas Melon, mesa mensa ‘Mensal’ Mines minged Midge, Medials, media? Minges mingle.","Linos/song songlike//naled naleds///name/ Lions nife//nifes///solid","","","Medal Nodal, node/. Nodi medial//Median. Muled Muon/na, Nail. nailed medina. nod Media, nodule, noel, noil//nole nomad/Melano melanoid meld/."}', 
                    '{"images/lessons/bottom_row_map.svg","","","","","","",""}')
                ;

            INSERT INTO lessons (lessonname, lessontext, lessonid, chapterid, minimumscoretopass, lessondescriptions, image)
            VALUES
                ('Lesson 4: Left and Right Hand', 
                    '{"We will now combine your left and right hands for the bottom row. Remember to to return your fingers back to the home row when you are not using them!","Type the words below.","Awesome! Let’s try one more practice with the bottom row."}', 
                    uuid_generate_v4(), 
                    chapter_4, 
                    '{0, 90, 90}', 
                    '{"","Conium. Zinc, camion//cabmen. nix nib Nab mux Mob bucine Zobu zine Cinema “Abune” Above&#59 zone? bean Beam. bani/Zoic zoea.","Clanks As an Bam chalks. Am cam, Chanks, zax ash? Alb, abs Za Blanks. blanch Blacks zacks/. Zas van Sax ax//Cab ban?"}', 
                    '{"","",""}')
                ;

            INSERT INTO lessons (lessonname, lessontext, lessonid, chapterid, minimumscoretopass, lessondescriptions, image)
            VALUES
                ('Chapter 4 Test', 
                    '{"Now that you have learned all of the alphabet keys on the keyboard, this chapter test will review the use of  all of them!","Hopefully by now, you feel very comfortable with touch typing! Continue to improve your speed, accuracy, and techniques for this chapter test and try not to look down at the keyboard when completing the exercises!","This exercise contains long words! Good luck!","Now we will focus on shorter words. Try to find a good rhythm!","Type the words below!"}', 
                    uuid_generate_v4(), 
                    chapter_4, 
                    '{0,0,90,90,90}', 
                    '{"","","Bowstringed Whitecombs “whiteboard” Bodysurfing ‘Formatively’ forjudgment/? Abridgments Decryptions! Elucidators earth/moving Zymographs zygopterid, abolishment, campgrounds, chimneypots decurvation // boulderings.","Bald chis, Lade/lacy chip “Chin” Able, mink, Mine. Mind flam fax max Gap’s Gapo “Ions” into Gape! gaol: Bake Mils&#59 milk bait Hawk have abid//Abet lads Lack/. zoa Zit bail.","Falx Balm, Flab, flax. calx Aeon “Tony” Wait naoi Fan pal nap Lap alp Van can Lac//ban lab nab Cab Lam equal Panel, uveal, plane Mac, mal van. bal lax, man lav Zax! bac vac bam clam Calm. lamb? calf clash/of/clan!"}', 
                    '{"","","","",""}')
                ;


            -- Chapter 5 --
            INSERT INTO lessons (lessonname, lessontext, lessonid, chapterid, minimumscoretopass, lessondescriptions, image)
            VALUES
                ('Lesson 1: Introduction', 
                    '{"Now that you have mastered all of the alphabet keys on the keyboard, now it is time to learn how to type numbers on the top of your keyboard!","This chapter is designed to use the number row, the number keys on the top of your keyboard, not the numeric pad.","Make sure to move only one finger at a time! The other fingers should be resting on the home row at all times."}', 
                    uuid_generate_v4(), 
                    chapter_5, 
                    '{0,0,0}', 
                    '{"","",""}', 
                    '{"","images/lessons/number_row_map.svg",""}')
                ;

            INSERT INTO lessons (lessonname, lessontext, lessonid, chapterid, minimumscoretopass, lessondescriptions, image)
            VALUES
                ('Lesson 2: Left Hand (12345)', 
                    '{"The numbers ‘1, 2, 3, 4, 5’ can be reached with your left hand. Please see the chart below to learn which fingers are used to type which keys. Use pinky for ‘1’, ring for ‘2’, middle for 3’, and pointer for ‘4’ and ‘5’.","Let’s start with number ‘1’!","Great! Let’s move onto ‘2’!","Awesome work! Let’s practice with ‘3’!","Great! Let’s finish with both ‘4’ and ‘5’!","Awesome work. Let’s finish with all of the numbers on the left hand side. Remember to return your fingers to the home row."}', 
                    uuid_generate_v4(), 
                    chapter_5, 
                    '{0,80,80,80,80,80}', 
                    '{"","111 111 aaa a11 1a1 aa1 11a aa111 1a1 1 1 1 a a 1 a 1 aa 1 1","2 sss 222 s 2 s 2 s 22 s s 2 2 s 2 ss 2 2 222 s 2 ss 2 2 s 22 s s s 22","3333 dd 3 d 3 ddd 3 d 3 d 3 333 d d 3 ddd 3 3 3 d d3 3 d 3 d d dd3 3 3","f5 555 f f f  5 f 5 fff 5f 5 5 5 f 5 g 4 f 4 f 444 g 5 4 f 5 4 4 g 44 5 f 4 5 f5 5 4f","1a1a 2s2s 3d3d 4f4f 5f5f 11q 22w 33e 44r 55r 11z 22x 33c 44v 55v 12345"}', 
                    '{"images/lessons/number_row_map.svg","","","","",""}')
                ;

            INSERT INTO lessons (lessonname, lessontext, lessonid, chapterid, minimumscoretopass, lessondescriptions, image)
            VALUES
                ('Lesson 3: Right Hand (67890)', 
                    '{"The numbers ‘6, 7, 8, 9, 0’ can be reached with your right hand. Please see the chart below to learn which fingers are used to type which keys. Use pinky for ‘0’, ring for ‘9’, middle for ‘8’, and pointer for ‘7’ and ‘6’.","Let’s start with number ‘0’!","Great! Let’s move onto ‘9’!","Awesome! Let’s practice with ‘8’.","Great work. Let’s finish with both ‘7’ and ‘6’.","Awesome work! Let’s finish with all of the numbers on the right hand side. Remember to return your fingers to the home row."}', 
                    uuid_generate_v4(), 
                    chapter_5, 
                    '{0,80,80,80,80,80}', 
                    '{"","000 &#59 0&#59&#59 000 &#59 0&#59&#59 &#59&#590 &#590&#59 0&#590 00&#59 0&#59&#59 &#590&#59 000&#59 00 0 0 0 0&#590","99l ll9 9l9 9l9 999 l9l9 l999 l9ll 9l9 l9 l9l999 l9l 9l9l l9l9l","88k8k8k 888k8kk8 888k8 8k8 88k8 88k8 8k888 k8k8kk 888","7jjj 7j7j 77j7 7j7 777 j6 6j66 66jj6j6j 6jj6jj66j6 6j6j6 6 6 6 6j6j6 j7j 7j 77j 7","6j6j 77j7j 8jk8k8k 9l9l9l 0&#590&#5909k8j7j6h6j7j 9j9 h8j 8j8j 0l 0l0l k9 k"}', 
                    '{"images/lessons/number_row_map.svg","","","","",""}')
                ;

            INSERT INTO lessons (lessonname, lessontext, lessonid, chapterid, minimumscoretopass, lessondescriptions, image)
            VALUES
                ('Lesson 4: Left and Right Hand', 
                    '{"Before we combine the left and right hands, don’t worry too much about hitting the keys on the number row with the designated fingers. Unlike the alphabet keys, you have more flexibility when you reach for the number keys!","Type the words below!","Great! Let’s try some more examples."}', 
                    uuid_generate_v4(), 
                    chapter_5, 
                    '{0,90,90}', 
                    '{"","6 times 6 is 36. 7 minus 2 is five. 10 times 10 is 100. 6 divided by 2 is 3. 2 plus 4 times 8 times 7 times 0 is 2. 8 plus 2 is 10. 99 plus 2 is 101.","24 divided by 8 is 3. 1 plus 1 is 2. 9 minus 9 is 0. 5 times 5 is 25. 6 times 6 is 36. 9 plus 1 is 10. 72 - 4 is 68. 0 times 100 is 0. 3 times 3 is 9. 7 minus 0 is 7."}', 
                    '{"","",""}')
                ;

            INSERT INTO lessons (lessonname, lessontext, lessonid, chapterid, minimumscoretopass, lessondescriptions, image)
            VALUES
                ('Lesson 5: Arithmetic Operations', 
                    '{"In this lesson, you will learn how to type using basic arithmetic operations. Arithmetic operations implies the study of numbers with traditional operations, such as addition, subtraction, multiplication and division.","We will be using (+) for plus, (-) for minus, (*) for multiplication, (/) for division, and (=) for equal sign. Please see the map below on how to reach for these keys.","Type the arithmetic statements below.","Great! Let’s try some more examples"}', 
                    uuid_generate_v4(), 
                    chapter_5, 
                    '{0,0,90,90}', 
                    '{"","","6 * 6 = 36, 6 / 6 = 1, 7 - 2 = 5, 10 * 10 = 100, 6 / 2 = 3, 2 + 4 * 8 * 7 * 0 = 2, 8 + 2 = 10","24 / 8 = 3, 1 + 1 = 2, 9 - 9 = 0, 5 * 5 = 25, 6 * 6 = 36, 9 + 1 = 10, 27 - 7 = 20, 0 * 100 = 0, 9 - 2 = 7"}', 
                    '{"","images/lessons/keyboard_layout_map.svg","",""}')
                ;

            INSERT INTO lessons (lessonname, lessontext, lessonid, chapterid, minimumscoretopass, lessondescriptions, image)
            VALUES
                ('Chapter 5 Test', 
                    '{"Type the letters and number below.","Type the letters and numbers below."}', 
                    uuid_generate_v4(), 
                    chapter_5, 
                    '{90,90}', 
                    '{"7 minus 2 is five. 10 times 10 is 100. 2 + 4 * 8 * 7 * 0 = 2, 8 + 2 = 10. 9 plus 1 is 10. 72 - 4 is 68. 24 / 8 = 3, 1 + 1 = 2, 9 - 9 = 0. 72 + 8 = 80","64 - 4 = 60. 64 minus 4 is 60. 29 - 2 = 27. 29 minus 2 is 27. 30 + 30 = 60. 30 plus 30 is 60. 12 + 21 = 33. 12 plus 21 is 33. 6 / 2 = 3. 6 divided 2 is 3. 1 * 1 = 1. 1 times 1 is 1."}', 
                    '{"",""}')
                ;


            -- Chapter 6 --
            INSERT INTO lessons (lessonname, lessontext, lessonid, chapterid, minimumscoretopass, lessondescriptions, image)
            VALUES
                ('Lesson 1: Introduction', 
                    '{"Congratulations on making it to the final chapter! This chapter contains lessons and exercises that combines everything that you have learned. By completing these lessons and exercises, you will be able to increase your WPM, accuracy, and overall typing skills."}', 
                    uuid_generate_v4(), 
                    chapter_6, 
                    '{0}',
                    '{""}', 
                    '{""}')
                ;

            INSERT INTO lessons (lessonname, lessontext, lessonid, chapterid, minimumscoretopass, lessondescriptions, image)
            VALUES
                ('Lesson 1: Common Words', 
                    '{"Practicing is all about repetition. That is why this lesson is focused on common English words that are used today. By practicing common words, you will learn to memorize them and type faster!","“The” and “be” are two of the most used words!","Now let’s practice more common words!","Great! Now let’s practice more common words.","Awesome work! Now’ let’s practice more common words.","Great job! We just covered the top 30 most used English words!","Type the commonly used words below."}', 
                    uuid_generate_v4(), 
                    chapter_6, 
                    '{0,90,90,90,90,0,95}', 
                    '{"","The the the the The The the the the The be Be Be be be Be The The the the the Be Be be Be be Be be Be Be be The The the the Be be Be be The the the Be be Be","To to to To to to To and And and and And and and and a a a A a A A in In in in in In in In in that That that that that That have have have Have Have have I I I I The The the the The be Be be Be be","It it it it It it for For for for For not not not not Not on On On on with With With with with With he He he He he he He as As as As as As you You You You you you do Do Do Do do do do At At at at At at at at At at at","This this this this this this This but But But But But But but his His his His his by By By By by From From from from they they They They they they we We We we we say Say Say her her her Her she She she She she She she She","","The the the Be Be be be To to to To and And and A in In in in hat That have have have I I I t It it for For not not Not on On On with With he He he He he he He as As you do Do Do Do at at At at this this This but But By By By by From From They they they we We We we Say Say her her her Her she"}', 
                    '{"","","","","","",""}')
                ;

            INSERT INTO lessons (lessonname, lessontext, lessonid, chapterid, minimumscoretopass, lessondescriptions, image)
            VALUES
                ('Lesson 2: Common Sentences', 
                    '{"In this lesson, you will be typing complete sentences. Most of our exercises until now have been words that are not related to each other. Now, we will practice frequently used sentences in English!","Sentences have a unique structure to them that you should become familiar with. For example, they always start with a capital letter and usually end with a period. By practicing with sentences, automatically reaching for the ‘Shift’ key at the beginning of each sentence will come naturally!","Greetings!","Planning meals!","Happy, happy, happy!"}', 
                    uuid_generate_v4(), 
                    chapter_6, 
                    '{0,0,90,90,90}', 
                    '{"","","Good morning! Good afternoon. Goodnight. How are you? How was your weekend? Not bad, thanks. I haven’t seen you in a while. I’m so pleased to meet you. I’m doing fine, thanks. How’s your week been? I’ll see you later!","Let’s grab lunch. I know a good place nearby. Let’s cook dinner tomorrow night. I’ll have the same. Let’s get a drink sometime. What do you want for lunch? What are they serving today? Let’s go to a restaurant. That grocery store sells fresh products.","I’m very happy right now. He is very happy. I feel great! This is so awesome! What would make you happy? She is so happy right now. When was your happiest moment? My goal in life is to make other people happy. I feel like a champion."}', 
                    '{"","","","",""}')
                ;

            INSERT INTO lessons (lessonname, lessontext, lessonid, chapterid, minimumscoretopass, lessondescriptions, image)
            VALUES
                ('Lesson 3: Long Passages', 
                    '{"This lesson will introduce you to longer passages. But don’t worry, the passages will cover fun facts about the ocean! So you will not only practice typing long passages, but also learn fun facts!","These passages will be longer than the ones you have seen and practiced earlier. Make sure to find a good rhythm while you type, so that you don’t get too tired halfway through!","Our Blue Planet.","Deepest known area of the world is...","What lives in the water?","The Pacific Ocean!","It is time to clean up our ocean!"}', 
                    uuid_generate_v4(), 
                    chapter_6, 
                    '{0,0,90,90,90,90,90}', 
                    '{"","","Three-fourths of the earth is covered with water. In fact, there is around 1,260,000,000,000,000,000,000 liters of water in our world. When astronauts first saw the planet from space, they could mostly see water, so they called it the ‘Blue Planet’.","The deepest known area of the earth’s oceans is known as the Mariana Trench. Its deepest point measures 11km. That’s a long dive down! Because it is so deep and difficult to travel to, more people have been to the moon than have explored the Mariana Trench!","While there are countless marine life forms known to man, there are many that have yet been discovered. Some scientists suggest that there could actually be millions of marine life forms out there.","The Pacific Ocean is the world’s largest body of water. Its area covers at least one-third of the total surface area of the earth. It’s so large that it covers more area than all land masses of the world combined! It is also home to some 25,000 islands.","There’s a huge “island” of trash floating around the northern area of the Pacific Ocean right now. It is called the Great Pacific Garbage patch. This plastic garbage island floats inside the center of the Pacific’s rotating ocean current."}', 
                    '{"","","","","","",""}')
                ;
    END;
$$ LANGUAGE plpgsql;

SELECT start();