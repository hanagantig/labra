CREATE DATABASE IF NOT EXISTS mego;
USE mego_api_test;

-- Пользователи (те, кто могут логиниться)
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    uuid CHAR(36) NOT NULL UNIQUE DEFAULT (UUID()),
    f_name VARCHAR(50),
    l_name VARCHAR(50),
    password VARCHAR(200) NOT NULL,
    birth_date DATE,
    gender ENUM('M', 'F') DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL
);

-- Пациенты
CREATE TABLE IF NOT EXISTS profiles (
    id INT AUTO_INCREMENT PRIMARY KEY,
    uuid CHAR(36) NOT NULL UNIQUE DEFAULT (UUID()),
    user_id INT DEFAULT NULL UNIQUE COMMENT 'Associated user uid',
    creator_user_id INT NOT NULL COMMENT 'Created by user uid',
    f_name VARCHAR(50),
    l_name VARCHAR(50),
    birth_date DATE,
    gender ENUM('M', 'F') DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL
);

-- Связь кто имеет доступ к какому пациенту
CREATE TABLE user_profiles (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    profile_id INT NOT NULL,
    label VARCHAR(50),
    access_level ENUM('owner', 'editor', 'viewer') NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL,
    UNIQUE (user_id, profile_id)
);

-- Контакты (email, phone), могут быть у users или profiles
CREATE TABLE contacts (
    id INT AUTO_INCREMENT PRIMARY KEY,
    type ENUM('email', 'phone') NOT NULL,
    value VARCHAR(50) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL
);

-- Связь контакта с users или profiles
CREATE TABLE linked_contacts (
    id INT AUTO_INCREMENT PRIMARY KEY,
    contact_id INT NOT NULL ,
    entity_type ENUM('user', 'profile') NOT NULL,
    entity_id VARCHAR(255) NOT NULL,
    is_primary BOOLEAN DEFAULT true,
    verified_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE (contact_id, entity_type)
);

CREATE TABLE IF NOT EXISTS labs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    lab_name VARCHAR(100) NOT NULL,
    address VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS marker_sets (
    id INT AUTO_INCREMENT PRIMARY KEY,
    set_name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS units (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    unit VARCHAR(50),
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS conversions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    from_unit_id INT NOT NULL,
    to_unit_id INT NOT NULL,
    rate DECIMAL(18, 10) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (from_unit_id) REFERENCES units(id),
    FOREIGN KEY (to_unit_id) REFERENCES units(id)
);

CREATE TABLE IF NOT EXISTS markers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    ref_range_min FLOAT,
    ref_range_max FLOAT,
    primary_color VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS set_markers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    set_id INT NOT NULL,
    marker_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (set_id) REFERENCES marker_sets(id),
    FOREIGN KEY (marker_id) REFERENCES markers(id)
);

CREATE TABLE IF NOT EXISTS checkups (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    profile_id INT NOT NULL,
    lab_id INT NOT NULL,
    status VARCHAR(10) NOT NULL,
    uploaded_file_id INT UNIQUE ,
    date DATE NOT NULL,
    comment TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS checkup_results (
    id INT AUTO_INCREMENT PRIMARY KEY,
    checkup_id INT NOT NULL,
    marker_id INT NOT NULL,
    undefined_marker VARCHAR(100),
    unit_id INT NOT NULL,
    undefined_unit VARCHAR(100),
    value FLOAT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (checkup_id) REFERENCES checkups(id)
);

CREATE TABLE `codes` (
     `id` int(11) NOT NULL AUTO_INCREMENT,
     `user_id` CHAR(36) NOT NULL,
     `object_type` varchar(50) NOT NULL,
     `object_id` varchar(255) NOT NULL,
     `code` varchar(255) NOT NULL,
     `used_at` TIMESTAMP,
     `expired_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     `deleted_at` TIMESTAMP,
     PRIMARY KEY (`id`),
    UNIQUE KEY (`object_type`, `object_id`, `code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin
  AUTO_INCREMENT=1;

USE mego_api_test;
USE mego;
CREATE TABLE sessions (
    id int(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
    token              VARCHAR(255) NOT NULL UNIQUE,
    user_uuid            CHAR(36) NOT NULL,
    session_id        CHAR(36) NOT NULL,
    device_id         CHAR(36) NOT NULL,
    replaced_by         VARCHAR(255),
    replaced_at         DATETIME,
    revoked_at         DATETIME,
    expires_at         DATETIME NOT NULL,
    created_at         DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE uploaded_files (
    id int(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id            INT NOT NULL,
    profile_id            INT NOT NULL,
    file_id VARCHAR(50) NOT NULL,
    pipeline_id VARCHAR(50) NOT NULL,
    fingerprint CHAR(64) NOT NULL UNIQUE,
    file_type VARCHAR(50) NOT NULL,
    source VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL,
    attempts_left INT NOT NULL DEFAULT 0,
    details TEXT,
    created_at         DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (1, 'Гемоглобин', null, null, '2025-02-08 20:15:49', '2025-06-08 14:40:38', 65280);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (2, 'Гематокрит', null, null, '2025-02-08 20:41:34', '2025-06-08 14:55:23', 255);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (3, 'Эритроциты', null, null, '2025-02-08 20:52:01', '2025-06-08 14:55:23', 33023);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (4, 'Лейкоциты', null, null, '2025-02-08 20:53:00', '2025-06-08 14:55:23', 16711680);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (5, 'Тромбоциты', null, null, '2025-03-04 19:58:35', '2025-06-14 21:53:52', 255);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (6, 'Средний объем эритроцитов', null, null, '2025-03-04 19:58:35', '2025-06-08 14:40:38', 65280);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (7, 'Среднее содержание гемоглобина в эритроците', null, null, '2025-03-04 19:58:35', '2025-06-08 14:55:23', 16776960);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (8, 'Средняя концентрация гемоглобина в эритроците', null, null, '2025-03-04 19:58:35', '2025-06-08 14:40:38', 65280);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (9, 'Ширина распределения эритроцитов', null, null, '2025-03-04 19:58:35', '2025-06-08 14:40:38', 65280);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (10, 'СОЭ', null, null, '2025-03-04 19:58:35', '2025-06-14 18:49:30', 16744703);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (11, 'Протромбиновое время', null, null, '2025-03-04 19:58:35', '2025-06-08 14:40:38', 65280);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (12, 'МНО', null, null, '2025-03-04 19:58:35', '2025-06-08 14:40:38', 65280);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (13, 'Фибриноген', null, null, '2025-03-04 19:58:35', '2025-06-08 14:40:38', 65280);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (14, 'D-димер', null, null, '2025-03-04 19:58:35', '2025-06-08 14:40:38', 65280);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (15, 'Глюкоза крови', null, null, '2025-03-04 19:58:35', '2025-06-08 14:40:38', 65280);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (16, 'Гликированный гемоглобин', null, null, '2025-03-04 19:58:35', '2025-06-08 14:40:38', 65280);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (17, 'Инсулин', null, null, '2025-03-04 19:58:35', '2025-06-08 14:40:38', 65280);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (18, 'С-пептид', null, null, '2025-03-04 19:58:35', '2025-06-08 14:40:38', 65280);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (19, 'Общий белок', null, null, '2025-03-04 19:58:35', '2025-06-08 14:40:38', 65280);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (20, 'Альбумин', null, null, '2025-03-04 19:58:35', '2025-06-08 14:40:38', 65280);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (21, 'Билирубин общий', null, null, '2025-03-04 19:58:35', '2025-06-08 14:40:38', 65280);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (22, 'Билирубин прямой', null, null, '2025-03-04 19:58:35', '2025-06-08 14:40:38', 65280);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (23, 'АЛТ', null, null, '2025-03-04 19:58:35', '2025-06-08 14:40:38', 65280);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (24, 'АСТ', null, null, '2025-03-04 19:58:35', '2025-06-08 14:40:38', 65280);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (25, 'Щелочная фосфатаза', null, null, '2025-03-04 19:58:35', '2025-06-08 14:40:38', 65280);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (26, 'Гамма-глутамилтрансфераза', null, null, '2025-03-04 19:58:35', '2025-06-08 14:40:38', 65280);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (27, 'Лактатдегидрогеназа', null, null, '2025-03-04 19:58:35', '2025-06-08 14:40:38', 65280);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (28, 'Мочевина', null, null, '2025-03-04 19:58:35', '2025-06-08 14:40:38', 65280);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (29, 'Креатинин', null, null, '2025-03-04 19:58:35', '2025-06-08 14:40:38', 65280);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (30, 'Липидный профиль', null, null, '2025-03-04 19:58:35', '2025-06-08 14:40:38', 65280);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (31, 'Железо сывороточное', null, null, '2025-03-04 19:58:35', '2025-06-08 14:40:38', 65280);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (32, 'Ферритин', null, null, '2025-03-04 19:58:35', '2025-06-08 14:40:38', 65280);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (33, 'ОЖСС', null, null, '2025-03-04 19:58:35', '2025-06-08 14:40:38', 65280);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (34, 'Кальций общий', null, null, '2025-03-04 19:58:35', '2025-06-08 14:40:38', 65280);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (35, 'Кальций ионизированный', null, null, '2025-03-04 19:58:35', '2025-06-08 14:40:38', 65280);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (36, 'Магний', null, null, '2025-03-04 19:58:35', '2025-06-08 14:40:38', 65280);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (37, 'Калий', null, null, '2025-03-04 19:58:35', '2025-06-08 14:40:38', 65280);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (38, 'Натрий', null, null, '2025-03-04 19:58:35', '2025-06-08 14:40:38', 65280);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (39, 'Хлориды', null, null, '2025-03-04 19:58:35', '2025-06-08 14:40:38', 65280);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (40, 'Фосфор', null, null, '2025-03-04 19:58:35', '2025-06-08 14:40:38', 65280);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (41, 'Витамин B12', null, null, '2025-03-04 19:58:35', '2025-06-08 14:40:38', 65280);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (42, 'Фолиевая кислота (B9)', null, null, '2025-03-04 19:58:35', '2025-06-08 14:40:38', 65280);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (43, 'Витамин D (25-ОН)', null, null, '2025-03-04 19:58:35', '2025-06-08 14:40:38', 65280);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (44, 'ТТГ (Тиреотропный гормон)', null, null, '2025-03-04 19:58:35', '2025-06-08 14:40:38', 65280);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (45, 'Т4 свободный', null, null, '2025-03-04 19:58:35', '2025-06-08 14:40:38', 65280);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (46, 'Т3 свободный', null, null, '2025-03-04 19:58:35', '2025-06-08 14:40:38', 65280);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (47, 'Анти-ТПО (антитела к тиреоидной пероксидазе)', null, null, '2025-03-04 19:58:35', '2025-06-08 14:40:38', 65280);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (48, 'Анти-ТГ (антитела к тиреоглобулину)', null, null, '2025-03-04 19:58:35', '2025-06-08 14:40:38', 65280);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (49, 'Среднее содержание Hb в эритроците', null, null, '2025-03-05 13:38:15', '2025-06-08 15:02:50', 65535);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (50, 'Средняя концентрация Hb в эритроците', null, null, '2025-03-05 13:38:30', '2025-06-08 14:40:38', 65280);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (51, 'Гетерогенность эритроцитов по объёму', null, null, '2025-03-05 13:38:52', '2025-06-08 14:40:38', 65280);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (52, 'Средний объём тромбоцитов', null, null, '2025-03-05 13:39:18', '2025-06-08 14:40:38', 65280);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (53, 'Гетерогенность тромбоцитов по объёму', null, null, '2025-03-05 13:39:32', '2025-06-08 14:40:38', 65280);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (54, 'Тромбокрит', null, null, '2025-03-05 13:39:52', '2025-06-08 14:40:38', 65280);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (55, 'Нейтрофилы', null, null, '2025-03-05 13:40:09', '2025-06-08 14:40:38', 65280);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (56, 'Эозинофилы', null, null, '2025-03-05 13:40:23', '2025-06-08 14:40:38', 65280);
INSERT INTO mego.markers (id, name, ref_range_min, ref_range_max, created_at, updated_at, primary_color) VALUES (57, 'Базофилы', null, null, '2025-03-05 13:40:34', '2025-06-08 14:40:38', 65280);

INSERT INTO mego.units (id, name, unit, created_at, updated_at) VALUES (11, '10^12/л', '10^12/л', '2025-02-09 08:34:05', '2025-02-09 08:34:05');
INSERT INTO mego.units (id, name, unit, created_at, updated_at) VALUES (12, 'г/л', 'г/л', '2025-03-05 13:50:17', '2025-03-05 13:50:17');
INSERT INTO mego.units (id, name, unit, created_at, updated_at) VALUES (13, 'фл', 'фл', '2025-03-05 13:50:35', '2025-03-05 13:50:35');
INSERT INTO mego.units (id, name, unit, created_at, updated_at) VALUES (14, 'пг', 'пг', '2025-03-05 13:50:44', '2025-03-05 13:50:44');
INSERT INTO mego.units (id, name, unit, created_at, updated_at) VALUES (15, '%', '%', '2025-03-05 13:50:56', '2025-03-05 13:50:56');
INSERT INTO mego.units (id, name, unit, created_at, updated_at) VALUES (16, '10^9/л', '10^9/л', '2025-03-05 13:51:10', '2025-03-05 13:51:10');
INSERT INTO mego.units (id, name, unit, created_at, updated_at) VALUES (17, 'мм/час', 'мм/час', '2025-03-05 13:51:32', '2025-03-05 13:51:32');
