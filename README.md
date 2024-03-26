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

etc/mysql/mariadb.conf.d/50-server.cnf 
[mysqld]
bind-address = 0.0.0.0

sudo systemctl restart mariadb


<!-- SQL -->

CREATE DATABASE wellkinson;
CREATE USER 'mark'@'%' IDENTIFIED BY '6853';
GRANT ALL PRIVILEGES ON wellkinson.* TO 'mark'@'%';
FLUSH PRIVILEGES;

CREATE TABLE `alarms` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uid` int(11) NOT NULL COMMENT '유저아이디',
  `type` tinyint(4) NOT NULL COMMENT '타입코드 1: 운동 2:약',
  `parent_id` int(11) DEFAULT NULL COMMENT '부모 pk',
  `body` text NOT NULL COMMENT '알람내용',
  `start_at` varchar(255) NOT NULL COMMENT '시작일',
  `end_at` varchar(255) NOT NULL COMMENT '종료알',
  `timestamp` varchar(255) NOT NULL COMMENT '알람시간',
  `week` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '알람요일' CHECK (json_valid(`week`)),
  `created` datetime DEFAULT current_timestamp() COMMENT '생성일',
  `updated` datetime DEFAULT current_timestamp() ON UPDATE current_timestamp() COMMENT '수정일',
  PRIMARY KEY (`id`),
  KEY `uid` (`uid`),
  CONSTRAINT `alarms_ibfk_1` FOREIGN KEY (`uid`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci COMMENT='사용자별 알람정보';

CREATE TABLE `app_versions` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `latest_version` varchar(40) NOT NULL,
  `android_link` varchar(255) DEFAULT NULL,
  `ios_link` varchar(255) DEFAULT NULL,
  `created` datetime DEFAULT current_timestamp(),
  `updated` datetime DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci COMMENT='연동된 이메일 정보';

CREATE TABLE `auth_codes` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `phone_number` varchar(40) NOT NULL,
  `code` varchar(40) NOT NULL,
  `created` datetime DEFAULT current_timestamp(),
  `updated` datetime DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci COMMENT='인증번호 정보';

CREATE TABLE `diet_presets` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uid` int(11) NOT NULL,
  `name` varchar(40) NOT NULL,
  `foods` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL CHECK (json_valid(`foods`)),
  `created` datetime DEFAULT current_timestamp(),
  `updated` datetime DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  KEY `uid` (`uid`),
  CONSTRAINT `diet_presets_ibfk_1` FOREIGN KEY (`uid`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

CREATE TABLE `diets` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uid` int(11) NOT NULL,
  `name` varchar(40) NOT NULL,
  `time` varchar(40) NOT NULL,
  `type` int(11) NOT NULL,
  `foods` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL CHECK (json_valid(`foods`)),
  `created` datetime DEFAULT current_timestamp(),
  `updated` datetime DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `date` varchar(40) NOT NULL COMMENT '식단 날짜',
  PRIMARY KEY (`id`),
  KEY `uid` (`uid`),
  CONSTRAINT `diets_ibfk_1` FOREIGN KEY (`uid`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci

CREATE TABLE `emotions` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uid` int(11) NOT NULL,
  `emotion` varchar(40) NOT NULL,
  `state` varchar(255) NOT NULL,
  `created` datetime DEFAULT current_timestamp(),
  `updated` datetime DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  KEY `uid` (`uid`),
  CONSTRAINT `emotions_ibfk_1` FOREIGN KEY (`uid`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

CREATE TABLE `exercises` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uid` int(11) NOT NULL,
  `title` varchar(40) NOT NULL,
  `exercise_end_at` varchar(40) NOT NULL,
  `exercise_start_at` varchar(40) NOT NULL,
  `plan_end_at` varchar(40) NOT NULL,
  `plan_start_at` varchar(40) NOT NULL,
  `use_alarm` tinyint(1) NOT NULL,
  `weekdays` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL CHECK (json_valid(`weekdays`)),
  `created` datetime DEFAULT current_timestamp(),
  `updated` datetime DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  KEY `uid` (`uid`),
  CONSTRAINT `exercises_ibfk_1` FOREIGN KEY (`uid`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;


CREATE TABLE `exercise_infos` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uid` int(11) NOT NULL,
  `date_performed` varchar(40) NOT NULL,
  `exercise_id` int(11) NOT NULL,
  `created` datetime DEFAULT current_timestamp(),
  `updated` datetime DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  KEY `uid` (`uid`),
  KEY `exercise_id` (`exercise_id`),
  CONSTRAINT `exercise_infos_ibfk_1` FOREIGN KEY (`uid`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `exercise_infos_ibfk_2` FOREIGN KEY (`exercise_id`) REFERENCES `exercises` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;


CREATE TABLE `face_exams` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `video_id` varchar(40) DEFAULT NULL COMMENT 'vimeo 비디오 id',
  `title` varchar(40) NOT NULL COMMENT '표정명',
  `type` tinyint(4) NOT NULL COMMENT '코드 1:기쁨 2:슬픔 3:놀람 4:분노',
  `created` datetime DEFAULT current_timestamp(),
  `updated` datetime DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci COMMENT='표정검사 베이스 테이블';

CREATE TABLE `face_exercises` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `video_id` varchar(40) DEFAULT NULL COMMENT 'vimeo 비디오 id',
  `title` varchar(40) NOT NULL COMMENT '동영상 제목',
  `type` tinyint(4) NOT NULL COMMENT '코드 1:기쁨 2:슬픔 3:놀람 4:분노',
  `created` datetime DEFAULT current_timestamp(),
  `updated` datetime DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci COMMENT='표정운동 베이스 테이블';

CREATE TABLE `face_scores` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uid` int(11) NOT NULL,
  `score` int(11) NOT NULL,
  `type` tinyint(4) NOT NULL,
  `created` datetime DEFAULT current_timestamp(),
  `updated` datetime DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  KEY `uid` (`uid`),
  CONSTRAINT `face_scores_ibfk_1` FOREIGN KEY (`uid`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

CREATE TABLE `images` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uid` int(11) NOT NULL COMMENT '유저 아이디',
  `url` varchar(255) NOT NULL COMMENT '이미지 url',
  `thumbnail_url` varchar(255) NOT NULL COMMENT '썸네일 url',
  `parent_id` int(11) NOT NULL COMMENT '부모 아이디',
  `type` tinyint(4) NOT NULL COMMENT '0:메인 프로필 1:식단',
  `level` tinyint(4) NOT NULL DEFAULT 0 COMMENT '상태 0:기본 10:삭제',
  `created` datetime DEFAULT current_timestamp() COMMENT '생성일',
  `updated` datetime DEFAULT current_timestamp() ON UPDATE current_timestamp() COMMENT '수정일',
  PRIMARY KEY (`id`),
  KEY `uid` (`uid`),
  CONSTRAINT `images_ibfk_1` FOREIGN KEY (`uid`) REFERENCES `users` (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci COMMENT='사용자별 업로드한 이미지 정보';

CREATE TABLE `inquires` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uid` int(11) NOT NULL,
  `email` varchar(40) NOT NULL,
  `title` varchar(40) NOT NULL,
  `content` text NOT NULL,
  `created` datetime DEFAULT current_timestamp(),
  `updated` datetime DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `level` tinyint(4) NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  KEY `uid` (`uid`),
  CONSTRAINT `inquires_ibfk_1` FOREIGN KEY (`uid`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

CREATE TABLE `inquire_replies` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uid` int(11) NOT NULL,
  `inquire_id` int(11) NOT NULL,
  `reply_type` tinyint(1) NOT NULL COMMENT '답변 or 추가 문의',
  `content` text NOT NULL,
  `created` datetime DEFAULT current_timestamp(),
  `updated` datetime DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `level` tinyint(4) NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  KEY `uid` (`uid`),
  KEY `inquire_id` (`inquire_id`),
  CONSTRAINT `inquire_replies_ibfk_1` FOREIGN KEY (`uid`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `inquire_replies_ibfk_2` FOREIGN KEY (`inquire_id`) REFERENCES `inquires` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;


CREATE TABLE `linked_emails` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `email` varchar(40) NOT NULL,
  `uid` int(11) NOT NULL,
  `created` datetime DEFAULT current_timestamp(),
  `updated` datetime DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `sns_type` tinyint(4) NOT NULL COMMENT '0:카카오 1:구글 2:애플',
  PRIMARY KEY (`id`),
  UNIQUE KEY `email` (`email`),
  KEY `uid` (`uid`),
  CONSTRAINT `linked_emails_ibfk_1` FOREIGN KEY (`uid`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci COMMENT='연동된 이메일 정보';

CREATE TABLE `main_services` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `title` varchar(40) NOT NULL COMMENT '서비스명',
  `level` tinyint(4) NOT NULL DEFAULT 0 COMMENT ' 0: 활성화 ,10:비활성화',
  `created` datetime DEFAULT current_timestamp(),
  `updated` datetime DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci COMMENT='이용하고 싶은 서비스 베이스';

CREATE TABLE `medicines` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uid` int(11) NOT NULL COMMENT '유저아이디',
  `timestamp` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '복용시간' CHECK (json_valid(`timestamp`)),
  `weekdays` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '알람요일' CHECK (json_valid(`weekdays`)),
  `dose` float NOT NULL COMMENT '복용량',
  `interval_type` tinyint(4) NOT NULL COMMENT '복용타입',
  `is_active` tinyint(1) NOT NULL COMMENT '활성화 여부',
  `least_store` float DEFAULT NULL COMMENT '최소 비축량',
  `medicine_type` varchar(40) NOT NULL COMMENT '약 타입',
  `name` varchar(40) NOT NULL COMMENT '약 이름',
  `store` float DEFAULT NULL COMMENT '비축량',
  `start_at` varchar(40) DEFAULT NULL COMMENT '시작일',
  `end_at` varchar(40) DEFAULT NULL COMMENT '종료일',
  `use_privacy` tinyint(1) NOT NULL COMMENT '개인정보 보호 알림 사용 여부',
  `created` datetime DEFAULT current_timestamp(),
  `updated` datetime DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  KEY `uid` (`uid`),
  CONSTRAINT `medicines_ibfk_1` FOREIGN KEY (`uid`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci COMMENT='약 정보';

CREATE TABLE `medicine_searches` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(40) NOT NULL COMMENT '약 이름',
  `created` datetime DEFAULT current_timestamp(),
  `updated` datetime DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci COMMENT='사용자별 약물 복용정보';

CREATE TABLE `medicine_takes` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uid` int(11) NOT NULL COMMENT '유저아이디',
  `date_taken` varchar(40) NOT NULL COMMENT '약 복용 일자',
  `time_taken` varchar(40) NOT NULL COMMENT '약 복용 시간',
  `dose` float NOT NULL COMMENT '복용량',
  `medicine_id` int(11) NOT NULL COMMENT '약물 pk',
  `created` datetime DEFAULT current_timestamp(),
  `updated` datetime DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  KEY `uid` (`uid`),
  KEY `medicine_id` (`medicine_id`),
  CONSTRAINT `medicine_takes_ibfk_1` FOREIGN KEY (`uid`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `medicine_takes_ibfk_2` FOREIGN KEY (`medicine_id`) REFERENCES `medicines` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci COMMENT='사용자별 약물 복용정보';

CREATE TABLE `notifications` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uid` int(11) NOT NULL,
  `type` tinyint(4) NOT NULL,
  `body` text NOT NULL,
  `is_read` tinyint(1) NOT NULL DEFAULT 0,
  `created` datetime DEFAULT current_timestamp(),
  `updated` datetime DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  KEY `uid` (`uid`),
  CONSTRAINT `notifications_ibfk_1` FOREIGN KEY (`uid`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

CREATE TABLE `sleep_alarms` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uid` int(11) NOT NULL COMMENT '유저아이디',
  `start_time` varchar(40) NOT NULL COMMENT '취침 시간',
  `end_time` varchar(40) NOT NULL COMMENT '기상 시간',
  `alarm_time` varchar(40) NOT NULL COMMENT '알람 시간',
  `is_active` tinyint(1) NOT NULL COMMENT '활성화 여부',
  `weekdays` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '알람 요일' CHECK (json_valid(`weekdays`)),
  `created` datetime DEFAULT current_timestamp(),
  `updated` datetime DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  KEY `uid` (`uid`),
  CONSTRAINT `sleep_alarms_ibfk_1` FOREIGN KEY (`uid`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci COMMENT='사용자별 수면알람 정보';

CREATE TABLE `sleep_times` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uid` int(11) NOT NULL COMMENT '유저아이디',
  `start_time` varchar(40) NOT NULL COMMENT '취침 시간',
  `end_time` varchar(40) NOT NULL COMMENT '기상 시간',
  `date_sleep` varchar(40) NOT NULL COMMENT '수면 일자',
  `created` datetime DEFAULT current_timestamp(),
  `updated` datetime DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  KEY `uid` (`uid`),
  CONSTRAINT `sleep_times_ibfk_1` FOREIGN KEY (`uid`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci COMMENT='사용자별 수면시간 정보';

CREATE TABLE `users` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `is_admin` tinyint(1) NOT NULL DEFAULT 0,
  `birthday` varchar(40) NOT NULL,
  `device_id` varchar(40) NOT NULL,
  `gender` tinyint(1) NOT NULL,
  `fcm_token` varchar(255) NOT NULL,
  `is_first` tinyint(1) NOT NULL DEFAULT 1,
  `name` varchar(40) NOT NULL,
  `phone_num` varchar(40) NOT NULL COMMENT '휴대폰번호',
  `use_auto_login` tinyint(1) NOT NULL DEFAULT 0,
  `use_privacy_protection` tinyint(1) NOT NULL DEFAULT 0,
  `use_sleep_tracking` tinyint(1) NOT NULL DEFAULT 0,
  `user_type` tinyint(4) NOT NULL,
  `email` varchar(40) NOT NULL COMMENT '이메일',
  `created` datetime DEFAULT current_timestamp(),
  `updated` datetime DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `sns_type` tinyint(4) NOT NULL COMMENT '0:카카오 1:구글 2:애플',
  PRIMARY KEY (`id`),
  UNIQUE KEY `phone_num` (`phone_num`),
  UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

CREATE TABLE `user_services` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uid` int(11) NOT NULL COMMENT '유저아이디',
  `service_id` int(11) NOT NULL COMMENT '서비스 아이디',
  `title` varchar(40) NOT NULL COMMENT '서비스명',
  `created` datetime DEFAULT current_timestamp(),
  `updated` datetime DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  KEY `uid` (`uid`),
  KEY `service_id` (`service_id`),
  CONSTRAINT `user_services_ibfk_1` FOREIGN KEY (`uid`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `user_services_ibfk_2` FOREIGN KEY (`service_id`) REFERENCES `main_services` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci COMMENT='유저별 이용하고 싶은 서비스';

CREATE TABLE `verified_numbers` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `phone_number` varchar(40) NOT NULL,
  `created` datetime DEFAULT current_timestamp(),
  `updated` datetime DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci COMMENT='인증완료 휴대폰 번호 정보';

CREATE TABLE `videos` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `duration` int(11) DEFAULT NULL COMMENT '동영상길이',
  `name` varchar(255) NOT NULL COMMENT '제목',
  `project_name` varchar(255) NOT NULL COMMENT '상위 폴더명',
  `project_id` varchar(255) DEFAULT NULL COMMENT '상위 폴더 id',
  `video_id` varchar(40) NOT NULL COMMENT 'vimeo 비디오 id',
  `thumbnail_url` varchar(255) DEFAULT NULL COMMENT '썸네일url',
  `created` datetime DEFAULT current_timestamp(),
  `updated` datetime DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci COMMENT='동영상 정보';

CREATE TABLE `vocal_scores` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uid` int(11) NOT NULL COMMENT '유저아이디',
  `score` int(11) NOT NULL COMMENT '점수',
  `type` tinyint(4) NOT NULL COMMENT '코드 1:a 2:e 3:i 4:o 5:u',
  `created` datetime DEFAULT current_timestamp(),
  `updated` datetime DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  KEY `uid` (`uid`),
  CONSTRAINT `vocal_scores_ibfk_1` FOREIGN KEY (`uid`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci COMMENT='사용자별 발성검사 점수';

CREATE TABLE `vocal_words` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `title` varchar(40) NOT NULL COMMENT '단어',
  `type` tinyint(4) NOT NULL COMMENT '코드 1:a 2:e 3:i 4:o 5:u',
  `created` datetime DEFAULT current_timestamp(),
  `updated` datetime DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci COMMENT='발성운동 단어 베이스';
