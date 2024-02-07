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