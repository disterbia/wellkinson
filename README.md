CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    is_admin BOOLEAN NOT NULL DEFAULT false COMMENT '관리자 여부',
    birthday VARCHAR(40) NOT NULL COMMENT '생일',
    device_id VARCHAR(40)  NOT NULL  COMMENT '기기아이디',
    gender BOOLEAN NOT NULL COMMENT '성별 1:남자 0:여자',
    fcm_token VARCHAR(255) NOT NULL COMMENT 'fcm토큰',
    is_first BOOLEAN NOT NULL DEFAULT true COMMENT '',
    name VARCHAR(40) NOT NULL COMMENT '이름',
    phone_num VARCHAR(40) NOT NULL COMMENT '휴댜폰번호',
    use_auto_login BOOLEAN NOT NULL DEFAULT false COMMENT '자동로그인 여부', 
    use_privacy_protection BOOLEAN NOT NULL DEFAULT false COMMENT '개인정보 보호 알림 사용 여부',
    use_sleep_tracking BOOLEAN NOT NULL DEFAULT false COMMENT '수면 트래킹 기능 사용 여부',
    user_type VARCHAR(40) NOT NULL COMMENT '사용자 타입',
    email VARCHAR(40) NOT NULL COMMENT '이메일',
    created DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '생성일',
    updated DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '수정일'
)engine=InnoDB default charset utf8 COMMENT = '사용자 정보';

CREATE TABLE alarms (
    id INT AUTO_INCREMENT PRIMARY KEY,
    uid INT NOT NULL COMMENT '유저아이디',
    type TINYINT NOT NULL COMMENT '타입코드 1: 운동 2:약',
    parent_id INT COMMENT '부모 pk',
    body TEXT NOT NULL COMMENT '알람내용',
	start_at VARCHAR(255) NOT NULL COMMENT '시작일',
    end_at VARCHAR(255) NOT NULL COMMENT '종료알',
    timestamp VARCHAR(255) NOT NULL COMMENT '알람시간',
    week JSON NOT NULL COMMENT '알람요일',
	created DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '생성일',
    updated DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '수정일',
    FOREIGN KEY (uid) REFERENCES users(id) ON DELETE CASCADE
)engine=InnoDB default charset utf8 COMMENT = '사용자별 알람정보';


CREATE TABLE notifications (
    id INT AUTO_INCREMENT PRIMARY KEY,
    uid INT NOT NULL  COMMENT '유저아이디',
    type VARCHAR(40)  NOT NULL  COMMENT '알람타입',
    body TEXT  NOT NULL  COMMENT '알람내용',
    is_read BOOLEAN NOT NULL DEFAULT false  COMMENT '확인여부',
	created DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '생성일',
    updated DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '수정일',
    FOREIGN KEY (uid) REFERENCES users(id) ON DELETE CASCADE
)engine=InnoDB default charset utf8 COMMENT= '발송한 푸쉬알림';

CREATE TABLE inquires(
    id INT AUTO_INCREMENT PRIMARY KEY,
    uid INT NOT NULL COMMENT '유저아이디',
    email VARCHAR(40) NOT NULL COMMENT '이메일',
    title VARCHAR(40)NOT NULL COMMENT '제목',
    level TINYINT NOT NULL DEFAULT 0 COMMENT '상태 0:기본 10:삭제',
    content TEXT NOT NULL COMMENT '내용',
    created DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '생성일',
    updated DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '수정일',
    FOREIGN KEY (uid) REFERENCES users(id) ON DELETE CASCADE
)engine=InnoDB default charset utf8 COMMENT = '사용자별 문의정보 ';

CREATE TABLE inquire_replies(
    id INT AUTO_INCREMENT PRIMARY KEY,
    uid INT NOT NULL COMMENT  '유저 아이디', 
    inquire_id INT NOT NULL COMMENT  '문의 pk',
    level TINYINT NOT NULL DEFAULT 0 COMMENT '상태 0:기본 10:삭제',
    reply_type BOOLEAN NOT NULL COMMENT '답변:1 추가문의:0',
    content TEXT NOT NULL COMMENT '내용',
    created DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '생성일',
    updated DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '수정일',
    FOREIGN KEY (uid) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (inquire_id) REFERENCES Inquires(id) ON DELETE CASCADE
)engine=InnoDB default charset utf8 COMMENT ='문의의 답변/추가문의';

CREATE TABLE diet_presets(
    id INT AUTO_INCREMENT PRIMARY KEY,
    uid INT NOT NULL COMMENT = '유저 아이디', 
    name VARCHAR(40) NOT NULL COMMENT = '식단명',
    foods json NOT NULL COMMENT = '음식 배열',
    created DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '생성일',
    updated DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '수정일',
    FOREIGN KEY (uid) REFERENCES users(id) ON DELETE CASCADE
)engine=InnoDB default charset utf8 COMMENT = '사용자별 식단계획 정보';

CREATE TABLE diets(
    id INT AUTO_INCREMENT PRIMARY KEY,
    uid INT NOT NULL COMMENT '유저 아이디', 
    name VARCHAR(40) NOT NULL COMMENT '식단명',
    time VARCHAR(40) NOT NULL COMMENT '식단 시간',
    type int NOT NULL COMMENT '아침/점심/저녁/간식 1/2/3/4' , 
    foods json NOT NULL COMMENT '음식 배열',
    created DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '생성일',
    updated DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '수정일',
    FOREIGN KEY (uid) REFERENCES users(id) ON DELETE CASCADE
)engine=InnoDB default charset utf8 COMMENT = '사용자별 먹은 식단 정보';

CREATE TABLE images(
    id INT AUTO_INCREMENT PRIMARY KEY,
    uid INT NOT NULL COMMENT '유저 아이디', 
    url VARCHAR(255) NOT NULL COMMENT '이미지 url',
    thumbnail_url VARCHAR(255) NOT NULL COMMENT '썸네일 url',
    diet_id INT COMMENT '식단 아이디' ,
    level TINYINT NOT NULL DEFAULT 0 COMMENT '상태 0:기본 10:삭제',
    created DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '생성일',
    updated DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '수정일',
    FOREIGN KEY (uid) REFERENCES users(id),
    FOREIGN KEY (diet_id) REFERENCES diets(id) 
)engine=InnoDB default charset utf8  COMMENT = '사용자별 업로드한 이미지 정보';

CREATE TABLE emotions(
    id INT AUTO_INCREMENT PRIMARY KEY,
    uid INT NOT NULL COMMENT = '유저 아이디', 
    emotion VARCHAR(40) NOT NULL COMMENT '기분명',
    state VARCHAR(255) NOT NULL COMMENT '기분 내용',
    created DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '생성일',
    updated DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '수정일',
    FOREIGN KEY (uid) REFERENCES users(id) ON DELETE CASCADE
)engine=InnoDB default charset utf8 COMMENT = '사용자별 기분 정보';

CREATE TABLE exercises(
    id INT AUTO_INCREMENT PRIMARY KEY,
    uid INT NOT NULL COMMENT '유저아이디', 
    title VARCHAR(40) NOT NULL COMMENT '운동명' ,
    exercise_end_at VARCHAR(40) NOT NULL  COMMENT '운동 종료시간',
    exercise_start_at VARCHAR(40) NOT NULL  COMMENT '운동 시작시간',
    plan_end_at VARCHAR(40) NOT NULL  COMMENT '운동 종료일자',
    plan_start_at VARCHAR(40) NOT NULL  COMMENT '운동 시작일자',
    use_alarm BOOLEAN NOT NULL  COMMENT '알람사용여부',
    weekdays JSON NOT NULL  COMMENT '운동 요일',
    created DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (uid) REFERENCES users(id) ON DELETE CASCADE
)engine=InnoDB default charset utf8 COMMENT = '사용자별 운동계획 정보';

CREATE TABLE exercise_infos(
    id INT AUTO_INCREMENT PRIMARY KEY,
    uid INT NOT NULL  COMMENT '유저아이디',
    date_performed VARCHAR(40) NOT NULL  COMMENT '운동완료 일자',
    exercise_id INT NOT NULL  COMMENT '운동 pk',
    created DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (uid) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (exercise_id) REFERENCES exercises(id) ON DELETE CASCADE
)engine=InnoDB default charset utf8 COMMENT = '사용자별 운동 실행정보';

CREATE TABLE face_scores(
    id INT AUTO_INCREMENT PRIMARY KEY,
    uid INT NOT NULL  COMMENT '유저아이디',  
    score INT NOT NULL  COMMENT '점수',
    type TINYINT NOT NULL  COMMENT '코드 1:기쁨 2:슬픔 3:놀람 4:분노', 
    created DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (uid) REFERENCES users(id) ON DELETE CASCADE
)engine=InnoDB default charset utf8 COMMENT = '사용자별 표정검사 점수';

CREATE TABLE face_exams(
    id INT AUTO_INCREMENT PRIMARY KEY,
    video_id VARCHAR(40) COMMENT 'vimeo 비디오 id' ,
    title INT NOT NULL  COMMENT '표정명',
    type TINYINT NOT NULL  COMMENT '코드 1:기쁨 2:슬픔 3:놀람 4:분노', 
    created DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
)engine=InnoDB default charset utf8 COMMENT = '표정검사 베이스 테이블';

CREATE TABLE videos(
    id INT AUTO_INCREMENT PRIMARY KEY,
    duration INT COMMENT '동영상길이', 
    name VARCHAR(255) NOT NULL  COMMENT '제목',
    project_name VARCHAR(255) NOT NULL  COMMENT '상위 폴더명',
    project_id VARCHAR(255) COMMENT '상위 폴더 id',
    video_id VARCHAR(40) NOT NULL COMMENT 'vimeo 비디오 id' ,
    thumbnail_url VARCHAR(255) COMMENT '썸네일url',
    created DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
)engine=InnoDB default charset utf8 COMMENT = '동영상 정보';


CREATE TABLE medicines(
    id INT AUTO_INCREMENT PRIMARY KEY,
    uid INT NOT NULL  COMMENT '유저아이디',  
    timestamp JSON COMMENT '복용시간',
    weekdays JSON COMMENT '알람요일',
    dose FLOAT NOT NULL COMMENT '복용량',
    interval_type TINYINT NOT NULL COMMENT '복용타입',
    is_active BOOLEAN  NOT NULL COMMENT '활성화 여부',
    least_store FLOAT COMMENT '최소 비축량',
    medicine_type VARCHAR(40) NOT NULL COMMENT '약 타입',
    name VARCHAR(40) NOT NULL COMMENT '약 이름',
    store FLOAT COMMENT '비축량',
    start_at VARCHAR(40) COMMENT '시작일',
    end_at VARCHAR(40) COMMENT '종료일',
    use_privacy BOOLEAN NOT NULL COMMENT  '개인정보 보호 알림 사용 여부',
    created DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (uid) REFERENCES users(id) ON DELETE CASCADE
)engine=InnoDB default charset utf8 COMMENT = '사용자별 약물 정보';


CREATE TABLE medicine_takes(
    id INT AUTO_INCREMENT PRIMARY KEY,
    uid INT NOT NULL  COMMENT '유저아이디',
    date_taken VARCHAR(40) NOT NULL  COMMENT '약 복용 일자',
    time_taken VARCHAR(40)  NOT NULL  COMMENT '약 복용 시간',
    dose FLOAT NOT NULL COMMENT '복용량',
    medicine_id INT NOT NULL  COMMENT '약물 pk',
    created DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (uid) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (medicine_id) REFERENCES medicines(id) ON DELETE CASCADE
)engine=InnoDB default charset utf8 COMMENT = '사용자별 약물 복용정보';

CREATE TABLE medicine_searches(
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(40) NOT NULL COMMENT '약 이름',
    created DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
)engine=InnoDB default charset utf8 COMMENT = '약물 검색 데이터';

CREATE TABLE sleep_alarms(
    id INT AUTO_INCREMENT PRIMARY KEY,
    uid INT NOT NULL  COMMENT '유저아이디',
    start_time VARCHAR(40) NOT NULL COMMENT '취침 시간',
    end_time VARCHAR(40) NOT NULL COMMENT '기상 시간',
    alarm_time VARCHAR(40) NOT NULL COMMENT '알람 시간',
    is_active BOOLEAN NOT NULL  COMMENT '활성화 여부',
    weekdays JSON NOT NULL COMMENT '알람 요일',
    created DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (uid) REFERENCES users(id) ON DELETE CASCADE
)engine=InnoDB default charset utf8 COMMENT = '사용자별 수면알람 정보';


CREATE TABLE sleep_times(
    id INT AUTO_INCREMENT PRIMARY KEY,
    uid INT NOT NULL  COMMENT '유저아이디',
	start_time VARCHAR(40) NOT NULL COMMENT '취침 시간',
	end_time VARCHAR(40) NOT NULL COMMENT '기상 시간',
    date_sleep VARCHAR(40) NOT NULL  COMMENT '수면 일자',
    created DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (uid) REFERENCES users(id) ON DELETE CASCADE
)engine=InnoDB default charset utf8 COMMENT = '사용자별 수면시간 정보';

CREATE TABLE face_exercises(  
    id INT AUTO_INCREMENT PRIMARY KEY,
    video_id VARCHAR(40) COMMENT 'vimeo 비디오 id' ,
    title VARCHAR(40) NOT NULL  COMMENT '동영상 제목',
    type TINYINT NOT NULL  COMMENT '코드 1:기쁨 2:슬픔 3:놀람 4:분노', 
    created DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
)engine=InnoDB default charset utf8 COMMENT = '표정운동 베이스 테이블';



INSERT INTO medicine_searches (name) VALUES
('이지레보정50'),
('이지레보정75'),
('마도파정'),
('비유프로정'),
('비유피-4정10mg'),
('비유피-4정20mg'),
('비프로정20mg'),
('삼성프로피베린정'),
('엔피베린정20mg'),
('유로나정10mg'),
('유로베린정'),
('유로콘정'),
('유로프로정'),
('마도파정125'),
('유로픽스정'),
('유리베린정'),
('유베린정'),
('유프베린정'),
('유피베린정'),
('이연프로피베린정'),
('이연프로피베린정10mg'),
('큐어베린정'),
('포베린정20mg'),
('프로베린정'),
('마도파확산정125'),
('프로베정'),
('프로베펙스정20mg'),
('프로빈정20mg'),
('프로시톨정'),
('프로피정'),
('프롭정20mg'),
('휴피베린정'),
('리보트릴정'),
('환인클로나제팜정0.5mg'),
('마도파에이치비에스캡슐125'),
('퍼킨씨알정25'),
('퍼킨씨알정50-200mg'),
('퍼킨씨알정25-100mg'),
('명도파정25'),
('명도파정50'),
('트리도파정50'),
('이지레보정100'),
('트리도파정75'),
('트리도파정100'),
('트리도파정125'),
('트리도파정150'),
('트리도파정200'),
('시네메트정'),
('시네메트씨알정'),
('시네메트정25'),
('미라펙스서방정0.375mg'),
('미라펙스서방정0.75mg'),
('이지레보정125'),
('미라펙스서방정1.5mg'),
('미라펙스정0.125mg'),
('미라펙스정0.25mg'),
('미라펙스정0.5mg'),
('미라펙스정1.0mg'),
('리큅정0.25mg'),
('리큅정1mg'),
('리큅정2mg'),
('리큅정5mg'),
('리큅피디정2mg'),
('스타레보필름코팅정50'),
('리큅피디정4mg'),
('리큅피디정8mg'),
('파키놀정0.25mg'),
('파키놀정1mg'),
('파키놀정2mg'),
('파키놀정5mg'),
('파키놀피디정2mg'),
('파키놀피디정4mg'),
('파키놀피디정8mg'),
('팔로델정'),
('스타레보필름코팅정75'),
('뉴큅정0.25mg'),
('뉴큅정1mg'),
('뉴큅정2mg'),
('로피맥스정0.25mg'),
('로피맥스정1mg'),
('로피맥스정2mg'),
('로피맥스피디정2mg'),
('로피맥스피디정4mg'),
('로피맥스피디정8mg'),
('도파프로정0.25mg'),
('스타레보필름코팅정100'),
('도파프로정1.0mg'),
('도파프로정2.0mg'),
('오니롤정0.25mg'),
('오니롤정1mg'),
('오니롤정2mg'),
('프라펙솔정0.125mg'),
('프라펙솔정0.25mg'),
('프라펙솔정0.5mg'),
('프라펙솔정1mg'),
('프라펙솔서방정0.75mg'),
('스타레보필름코팅정125'),
('프라펙솔서방정0.375mg'),
('프라펙솔서방정1.5mg'),
('피디펙솔정0.125mg'),
('피디펙솔정0.25mg'),
('피디펙솔정0.5mg'),
('피디펙솔정1mg'),
('피디펙솔서방정0.375mg'),
('피디펙솔서방정1.5mg'),
('마오비정'),
('콤탄정200mg'),
('스타레보필름코팅정150'),
('아질렉트정'),
('아질렉트정0.5mg'),
('트리헥신정'),
('환인벤즈트로핀정'),
('환인벤즈트로핀정1mg'),
('명인벤즈트로핀메실산염정2mg'),
('명인벤즈트로핀메실산염정1mg'),
('프로이머정'),
('파로마정'),
('영프로마정'),
('스타레보필름코팅정200'),
('피케이멜즈정'),
('아만타정'),
('신일비사코딜정'),
('락시티브정'),
('둘코락스에스장용정'),
('넥스베린정20mg'),
('글로베린정'),
('대웅바이오프로피베린정20mg'),
('디트루베린정'),
('베아베린정');

INSERT INTO face_exercises (video_id, title, type) VALUES
('795204724', '뺨 올리기', 1),
('795204919', '윙크', 1),
('795203942', '눈썹올리기', 2),
('795203091', '눈깜빡감기', 2),
('795202256', '눈 천천히 감기', 2),
('795205255', '아래턱 움직이기', 2),
('795205026', '혀로 뺨 밀기', 2),
('795203640', '혀 입끝으로 보내기', 2),
('795202539', '눈썹모으기', 3),
('795203091', '눈깜빡감기', 3),
('795202256', '눈 천천히 감기', 3),
('795202446', '입둘레 부풀리기', 3),
('795203459', '입술깔대기만들기', 3),
('795204129', '입술모아올리기', 3),
('795202814', '입술 말아 넣기', 3),
('795202539', '눈썹모으기', 4),
('795202946', '입꼬리 내리기', 4);