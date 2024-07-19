/home/ubuntu/.ssh/authorized_keys // 퍼블릭키 변경

sudo apt update
sudo apt install -y mariadb-server
sudo systemctl start mariadb
sudo systemctl enable mariadb
sudo mysql_secure_installation
sudo mysql -u root -p
sudo ufw allow 3306/tcp

sudo apt install net-tools
netstat -nlpt

/etc/mysql/mariadb.conf.d/50-server.cnf 
[mysqld]
bind-address = 0.0.0.0

sudo systemctl restart mariadb

sudo apt-get update
sudo apt-get install docker.io
sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose


docker-compose up -d
docker-compose pull // 도커헙 새로푸쉬되었을때 가져오고 docker-compose up -d

//sudo 빼는법
sudo grep docker /etc/group // 존재하는 그룹이 있다면 아래 명령어 없다면 sudo groupadd docker
sudo usermod -aG docker $USER
newgrp docker

<!-- SQL -->

CREATE DATABASE wellkinson;
CREATE USER 'mark'@'%' IDENTIFIED BY '6853';
GRANT ALL PRIVILEGES ON wellkinson.* TO 'mark'@'%';
FLUSH PRIVILEGES;

-- wellkinson.app_versions definition
CREATE TABLE app_versions (
  id SERIAL PRIMARY KEY,
  latest_version VARCHAR(40) NOT NULL,
  android_link VARCHAR(255),
  ios_link VARCHAR(255),
  created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON TABLE app_versions IS '연동된 이메일 정보';

-- wellkinson.auth_codes definition
CREATE TABLE auth_codes (
  id SERIAL PRIMARY KEY,
  phone_number VARCHAR(40) NOT NULL,
  code VARCHAR(40) NOT NULL,
  created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON TABLE auth_codes IS '인증번호 정보';

-- wellkinson.face_exams definition
CREATE TABLE face_exams (
  id SERIAL PRIMARY KEY,
  video_id VARCHAR(40),
  title VARCHAR(40) NOT NULL,
  type SMALLINT NOT NULL,
  created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON TABLE face_exams IS '표정검사 베이스 테이블';
COMMENT ON COLUMN face_exams.video_id IS 'vimeo 비디오 id';
COMMENT ON COLUMN face_exams.title IS '표정명';
COMMENT ON COLUMN face_exams.type IS '코드 1:기쁨 2:슬픔 3:놀람 4:분노';

-- wellkinson.face_exercises definition
CREATE TABLE face_exercises (
  id SERIAL PRIMARY KEY,
  video_id VARCHAR(40),
  title VARCHAR(40) NOT NULL,
  type SMALLINT NOT NULL,
  created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  guide_video_id VARCHAR(40)
);
COMMENT ON TABLE face_exercises IS '표정운동 베이스 테이블';
COMMENT ON COLUMN face_exercises.video_id IS 'vimeo 비디오 id';
COMMENT ON COLUMN face_exercises.title IS '동영상 제목';
COMMENT ON COLUMN face_exercises.type IS '코드 1:기쁨 2:슬픔 3:놀람 4:분노';
COMMENT ON COLUMN face_exercises.guide_video_id IS '가이드 비디오 아이디';

-- wellkinson.main_services definition
CREATE TABLE main_services (
  id SERIAL PRIMARY KEY,
  title VARCHAR(40) NOT NULL,
  level SMALLINT NOT NULL DEFAULT 0,
  created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON TABLE main_services IS '이용하고 싶은 서비스 베이스';
COMMENT ON COLUMN main_services.title IS '서비스명';
COMMENT ON COLUMN main_services.level IS '0: 활성화 ,10:비활성화';

-- wellkinson.medicine_searches definition
CREATE TABLE medicine_searches (
  id SERIAL PRIMARY KEY,
  name VARCHAR(40) NOT NULL,
  created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON TABLE medicine_searches IS '사용자별 약물 복용정보';
COMMENT ON COLUMN medicine_searches.name IS '약 이름';

-- wellkinson.users definition
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  is_admin BOOLEAN NOT NULL DEFAULT FALSE,
  birthday VARCHAR(40) NOT NULL,
  device_id VARCHAR(40) NOT NULL,
  gender BOOLEAN NOT NULL,
  fcm_token VARCHAR(255) NOT NULL,
  is_first BOOLEAN NOT NULL DEFAULT TRUE,
  name VARCHAR(40) NOT NULL,
  phone_num VARCHAR(40) NOT NULL UNIQUE,
  use_auto_login BOOLEAN NOT NULL DEFAULT FALSE,
  use_privacy_protection BOOLEAN NOT NULL DEFAULT FALSE,
  use_sleep_tracking BOOLEAN NOT NULL DEFAULT FALSE,
  user_type SMALLINT NOT NULL,
  email VARCHAR(40) NOT NULL UNIQUE,
  created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  sns_type SMALLINT NOT NULL,
  indemnification_clause BOOLEAN NOT NULL
);
COMMENT ON TABLE users IS '연동된 이메일 정보';
COMMENT ON COLUMN users.phone_num IS '휴대폰번호';
COMMENT ON COLUMN users.email IS '이메일';
COMMENT ON COLUMN users.sns_type IS '0:카카오 1:구글 2:애플';
COMMENT ON COLUMN users.indemnification_clause IS '면책조항 동의 여부';

-- wellkinson.verified_numbers definition
CREATE TABLE verified_numbers (
  id SERIAL PRIMARY KEY,
  phone_number VARCHAR(40) NOT NULL,
  created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON TABLE verified_numbers IS '인증완료 휴대폰 번호 정보';

-- wellkinson.videos definition
CREATE TABLE videos (
  id SERIAL PRIMARY KEY,
  duration INT,
  name VARCHAR(255) NOT NULL,
  project_name VARCHAR(255) NOT NULL,
  project_id VARCHAR(255),
  video_id VARCHAR(40) NOT NULL,
  thumbnail_url VARCHAR(255),
  created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON TABLE videos IS '동영상 정보';
COMMENT ON COLUMN videos.duration IS '동영상길이';
COMMENT ON COLUMN videos.name IS '제목';
COMMENT ON COLUMN videos.project_name IS '상위 폴더명';
COMMENT ON COLUMN videos.project_id IS '상위 폴더 id';
COMMENT ON COLUMN videos.video_id IS 'vimeo 비디오 id';
COMMENT ON COLUMN videos.thumbnail_url IS '썸네일url';

-- wellkinson.vocal_words definition
CREATE TABLE vocal_words (
  id SERIAL PRIMARY KEY,
  title VARCHAR(40) NOT NULL,
  type SMALLINT NOT NULL,
  created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON TABLE vocal_words IS '발성운동 단어 베이스';
COMMENT ON COLUMN vocal_words.title IS '단어';
COMMENT ON COLUMN vocal_words.type IS '코드 1:a 2:e 3:i 4:o 5:u';

-- wellkinson.alarms definition
CREATE TABLE alarms (
  id SERIAL PRIMARY KEY,
  uid INT NOT NULL,
  type SMALLINT NOT NULL,
  parent_id INT,
  body TEXT NOT NULL,
  start_at VARCHAR(255) NOT NULL,
  end_at VARCHAR(255) NOT NULL,
  timestamp VARCHAR(255) NOT NULL,
  week JSONB NOT NULL CHECK (jsonb_typeof(week) = 'array'),
  created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT alarms_uid_fk FOREIGN KEY (uid) REFERENCES users (id) ON DELETE CASCADE
);
COMMENT ON TABLE alarms IS '사용자별 알람정보';
COMMENT ON COLUMN alarms.uid IS '유저아이디';
COMMENT ON COLUMN alarms.type IS '타입코드 1: 운동 2:약';
COMMENT ON COLUMN alarms.parent_id IS '부모 pk';
COMMENT ON COLUMN alarms.body IS '알람내용';
COMMENT ON COLUMN alarms.start_at IS '시작일';
COMMENT ON COLUMN alarms.end_at IS '종료알';
COMMENT ON COLUMN alarms.timestamp IS '알람시간';
COMMENT ON COLUMN alarms.week IS '알람요일';

-- wellkinson.diet_presets definition
CREATE TABLE diet_presets (
  id SERIAL PRIMARY KEY,
  uid INT NOT NULL,
  name VARCHAR(40) NOT NULL,
  foods JSONB NOT NULL CHECK (jsonb_typeof(foods) = 'array'),
  created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT diet_presets_uid_fk FOREIGN KEY (uid) REFERENCES users (id) ON DELETE CASCADE
);
COMMENT ON TABLE diet_presets IS '사용자별 다이어트 프리셋';

-- wellkinson.diets definition
CREATE TABLE diets (
  id SERIAL PRIMARY KEY,
  uid INT NOT NULL,
  memo VARCHAR(255),
  time VARCHAR(40) NOT NULL,
  type INT NOT NULL,
  foods JSONB NOT NULL CHECK (jsonb_typeof(foods) = 'array'),
  created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  date VARCHAR(40) NOT NULL,
  CONSTRAINT diets_uid_fk FOREIGN KEY (uid) REFERENCES users (id) ON DELETE CASCADE
);
COMMENT ON TABLE diets IS '사용자별 식단 정보';
COMMENT ON COLUMN diets.memo IS '메모';
COMMENT ON COLUMN diets.date IS '식단 날짜';

-- wellkinson.emotions definition
CREATE TABLE emotions (
  id SERIAL PRIMARY KEY,
  uid INT NOT NULL,
  emotion INT NOT NULL,
  state VARCHAR(255) NOT NULL,
  created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT emotions_uid_fk FOREIGN KEY (uid) REFERENCES users (id) ON DELETE CASCADE
);
COMMENT ON TABLE emotions IS '사용자별 감정 정보';
COMMENT ON COLUMN emotions.emotion IS '감정코드';

-- wellkinson.exercises definition
CREATE TABLE exercises (
  id SERIAL PRIMARY KEY,
  uid INT NOT NULL,
  title VARCHAR(40) NOT NULL,
  exercise_end_at VARCHAR(40) NOT NULL,
  exercise_start_at VARCHAR(40) NOT NULL,
  plan_end_at VARCHAR(40) NOT NULL,
  plan_start_at VARCHAR(40) NOT NULL,
  use_alarm BOOLEAN NOT NULL,
  weekdays JSONB NOT NULL CHECK (jsonb_typeof(weekdays) = 'array'),
  created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  is_delete BOOLEAN NOT NULL DEFAULT FALSE,
  CONSTRAINT exercises_uid_fk FOREIGN KEY (uid) REFERENCES users (id) ON DELETE CASCADE
);
COMMENT ON TABLE exercises IS '사용자별 운동 정보';
COMMENT ON COLUMN exercises.is_delete IS '삭제여부';

-- wellkinson.face_scores definition
CREATE TABLE face_scores (
  id SERIAL PRIMARY KEY,
  uid INT NOT NULL,
  score INT NOT NULL,
  type SMALLINT NOT NULL,
  created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT face_scores_uid_fk FOREIGN KEY (uid) REFERENCES users (id) ON DELETE CASCADE
);
COMMENT ON TABLE face_scores IS '사용자별 표정 점수';

-- wellkinson.images definition
CREATE TABLE images (
  id SERIAL PRIMARY KEY,
  uid INT NOT NULL,
  url VARCHAR(255) NOT NULL,
  thumbnail_url VARCHAR(255) NOT NULL,
  parent_id INT NOT NULL,
  type SMALLINT NOT NULL,
  level SMALLINT NOT NULL DEFAULT 0,
  created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT images_uid_fk FOREIGN KEY (uid) REFERENCES users (id)
);
COMMENT ON TABLE images IS '사용자별 업로드한 이미지 정보';
COMMENT ON COLUMN images.uid IS '유저 아이디';
COMMENT ON COLUMN images.url IS '이미지 url';
COMMENT ON COLUMN images.thumbnail_url IS '썸네일 url';
COMMENT ON COLUMN images.parent_id IS '부모 아이디';
COMMENT ON COLUMN images.type IS '0:메인 프로필 1:식단';
COMMENT ON COLUMN images.level IS '상태 0:기본 10:삭제';

-- wellkinson.inquires definition
CREATE TABLE inquires (
  id SERIAL PRIMARY KEY,
  uid INT NOT NULL,
  email VARCHAR(40) NOT NULL,
  title VARCHAR(40) NOT NULL,
  content TEXT NOT NULL,
  created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  level SMALLINT NOT NULL DEFAULT 0,
  CONSTRAINT inquires_uid_fk FOREIGN KEY (uid) REFERENCES users (id) ON DELETE CASCADE
);
COMMENT ON TABLE inquires IS '사용자별 문의 정보';

-- wellkinson.linked_emails definition
CREATE TABLE linked_emails (
  id SERIAL PRIMARY KEY,
  email VARCHAR(40) NOT NULL,
  uid INT NOT NULL,
  created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  sns_type SMALLINT NOT NULL,
  CONSTRAINT linked_emails_uid_fk FOREIGN KEY (uid) REFERENCES users (id) ON DELETE CASCADE,
  UNIQUE (email)
);
COMMENT ON TABLE linked_emails IS '연동된 이메일 정보';
COMMENT ON COLUMN linked_emails.sns_type IS '0:카카오 1:구글 2:애플';

-- wellkinson.medicines definition
CREATE TABLE medicines (
  id SERIAL PRIMARY KEY,
  uid INT NOT NULL,
  timestamp JSONB CHECK (jsonb_typeof(timestamp) = 'array'),
  weekdays JSONB CHECK (jsonb_typeof(weekdays) = 'array'),
  dose FLOAT NOT NULL,
  interval_type SMALLINT NOT NULL,
  is_active BOOLEAN NOT NULL,
  least_store FLOAT,
  medicine_type VARCHAR(40) NOT NULL,
  name VARCHAR(40) NOT NULL,
  store FLOAT,
  start_at VARCHAR(40),
  end_at VARCHAR(40),
  use_privacy BOOLEAN NOT NULL,
  created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  use_least_store BOOLEAN NOT NULL,
  is_delete BOOLEAN NOT NULL DEFAULT FALSE,
  CONSTRAINT medicines_uid_fk FOREIGN KEY (uid) REFERENCES users (id) ON DELETE CASCADE
);
COMMENT ON TABLE medicines IS '약 정보';

-- wellkinson.notifications definition
CREATE TABLE notifications (
  id SERIAL PRIMARY KEY,
  uid INT NOT NULL,
  type SMALLINT NOT NULL,
  body TEXT NOT NULL,
  is_read BOOLEAN NOT NULL DEFAULT FALSE,
  created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  parent_id INT,
  CONSTRAINT notifications_uid_fk FOREIGN KEY (uid) REFERENCES users (id) ON DELETE CASCADE
);
COMMENT ON TABLE notifications IS '알림 정보';

-- wellkinson.sleep_alarms definition
CREATE TABLE sleep_alarms (
  id SERIAL PRIMARY KEY,
  uid INT NOT NULL,
  start_time VARCHAR(40) NOT NULL,
  end_time VARCHAR(40) NOT NULL,
  alarm_time VARCHAR(40) NOT NULL,
  is_active BOOLEAN NOT NULL,
  weekdays JSONB NOT NULL CHECK (jsonb_typeof(weekdays) = 'array'),
  created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT sleep_alarms_uid_fk FOREIGN KEY (uid) REFERENCES users (id) ON DELETE CASCADE
);
COMMENT ON TABLE sleep_alarms IS '사용자별 수면알람 정보';

-- wellkinson.sleep_times definition
CREATE TABLE sleep_times (
  id SERIAL PRIMARY KEY,
  uid INT NOT NULL,
  start_time VARCHAR(40) NOT NULL,
  end_time VARCHAR(40) NOT NULL,
  date_sleep VARCHAR(40) NOT NULL,
  created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT sleep_times_uid_fk FOREIGN KEY (uid) REFERENCES users (id) ON DELETE CASCADE
);
COMMENT ON TABLE sleep_times IS '사용자별 수면시간 정보';

-- wellkinson.user_services definition
CREATE TABLE user_services (
  id SERIAL PRIMARY KEY,
  uid INT NOT NULL,
  service_id INT NOT NULL,
  title VARCHAR(40) NOT NULL,
  created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT user_services_uid_fk FOREIGN KEY (uid) REFERENCES users (id) ON DELETE CASCADE,
  CONSTRAINT user_services_service_id_fk FOREIGN KEY (service_id) REFERENCES main_services (id) ON DELETE CASCADE
);
COMMENT ON TABLE user_services IS '유저별 이용하고 싶은 서비스';

-- wellkinson.vocal_scores definition
CREATE TABLE vocal_scores (
  id SERIAL PRIMARY KEY,
  uid INT NOT NULL,
  score INT NOT NULL,
  type SMALLINT NOT NULL,
  created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT vocal_scores_uid_fk FOREIGN KEY (uid) REFERENCES users (id) ON DELETE CASCADE
);
COMMENT ON TABLE vocal_scores IS '사용자별 발성검사 점수';

-- wellkinson.exercise_infos definition
CREATE TABLE exercise_infos (
  id SERIAL PRIMARY KEY,
  uid INT NOT NULL,
  date_performed VARCHAR(40) NOT NULL,
  exercise_id INT NOT NULL,
  created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT exercise_infos_uid_fk FOREIGN KEY (uid) REFERENCES users (id) ON DELETE CASCADE,
  CONSTRAINT exercise_infos_exercise_id_fk FOREIGN KEY (exercise_id) REFERENCES exercises (id) ON DELETE CASCADE
);
COMMENT ON TABLE exercise_infos IS '사용자별 운동 정보';

-- wellkinson.inquire_replies definition
CREATE TABLE inquire_replies (
  id SERIAL PRIMARY KEY,
  uid INT NOT NULL,
  inquire_id INT NOT NULL,
  reply_type BOOLEAN NOT NULL,
  content TEXT NOT NULL,
  created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  level SMALLINT NOT NULL DEFAULT 0,
  CONSTRAINT inquire_replies_uid_fk FOREIGN KEY (uid) REFERENCES users (id) ON DELETE CASCADE,
  CONSTRAINT inquire_replies_inquire_id_fk FOREIGN KEY (inquire_id) REFERENCES inquires (id) ON DELETE CASCADE
);
COMMENT ON TABLE inquire_replies IS '사용자별 문의 답변 정보';

-- wellkinson.medicine_takes definition
CREATE TABLE medicine_takes (
  id SERIAL PRIMARY KEY,
  uid INT NOT NULL,
  date_taken VARCHAR(40) NOT NULL,
  time_taken VARCHAR(40) NOT NULL,
  dose FLOAT NOT NULL,
  medicine_id INT NOT NULL,
  created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  real_taken VARCHAR(40) NOT NULL,
  CONSTRAINT medicine_takes_uid_fk FOREIGN KEY (uid) REFERENCES users (id) ON DELETE CASCADE,
  CONSTRAINT medicine_takes_medicine_id_fk FOREIGN KEY (medicine_id) REFERENCES medicines (id) ON DELETE CASCADE
);
COMMENT ON TABLE medicine_takes IS '사용자별 약물 복용 정보';
