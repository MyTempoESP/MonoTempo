-- phpMyAdmin SQL Dump
-- version 5.2.1
-- https://www.phpmyadmin.net/
--
-- Host: db
-- Generation Time: Oct 08, 2024 at 07:39 PM
-- Server version: 8.0.39
-- PHP Version: 8.2.8

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `mytempo`
--

-- --------------------------------------------------------

--
-- Table structure for table `athletes`
--

CREATE TABLE `athletes` (
  `num` int NOT NULL,
  `event_id` int DEFAULT NULL,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `city` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `team` varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `track_id` int DEFAULT NULL,
  `sex` varchar(1) COLLATE utf8mb4_unicode_ci NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------

--
-- Table structure for table `athletes_times`
--

CREATE TABLE `athletes_times` (
  `id` int NOT NULL,
  `antenna` int DEFAULT NULL,
  `checkpoint_id` int NOT NULL,
  `athlete_num` int NOT NULL,
  `athlete_time` varchar(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `staff` int NOT NULL,
  `timestp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------

--
-- Table structure for table `checkpoints`
--

CREATE TABLE `checkpoints` (
  `id` int NOT NULL,
  `description` int NOT NULL,
  `km` int NOT NULL,
  `local` int NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- --------------------------------------------------------

--
-- Table structure for table `classificacao`
--

CREATE TABLE `classificacao` (
  `athlete` int DEFAULT NULL,
  `start_time` varchar(12) DEFAULT NULL,
  `end_time` varchar(12) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- --------------------------------------------------------

--
-- Table structure for table `equipamento`
--

CREATE TABLE `equipamento` (
  `id` tinyint(1) NOT NULL,
  `checkpoint_id` int NOT NULL,
  `idequip` int NOT NULL,
  `modelo` varchar(30) NOT NULL,
  `event_id` int NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- --------------------------------------------------------

--
-- Table structure for table `event_data`
--

CREATE TABLE `event_data` (
  `id` int NOT NULL,
  `event_date` date DEFAULT NULL,
  `event_title` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------

--
-- Table structure for table `invalidos`
--

CREATE TABLE `invalidos` (
  `id` int NOT NULL,
  `antenna` int DEFAULT NULL,
  `checkpoint_id` int NOT NULL,
  `athlete_num` int NOT NULL,
  `athlete_time` varchar(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `staff` int NOT NULL,
  `timestp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------

--
-- Table structure for table `rede`
--

CREATE TABLE `rede` (
  `ssid` varchar(100) NOT NULL,
  `id` int NOT NULL,
  `password` int NOT NULL,
  `status` tinyint NOT NULL DEFAULT '0',
  `descricao` varchar(255) NOT NULL,
  `conectar` tinyint NOT NULL DEFAULT '0'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- --------------------------------------------------------

--
-- Table structure for table `staffs`
--

CREATE TABLE `staffs` (
  `id` int NOT NULL,
  `event_id` int NOT NULL,
  `nome` varchar(200) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- --------------------------------------------------------

--
-- Table structure for table `tracks`
--

CREATE TABLE `tracks` (
  `id` int NOT NULL,
  `event_id` int DEFAULT NULL,
  `race_description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `inicio` time DEFAULT NULL,
  `chegada` time DEFAULT NULL,
  `largada` time DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------

--
-- Table structure for table `unsent_athletes`
--

CREATE TABLE `unsent_athletes` (
  `id` int NOT NULL,
  `antenna` int DEFAULT NULL,
  `checkpoint_id` int NOT NULL,
  `athlete_num` int NOT NULL,
  `athlete_time` varchar(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `staff` int NOT NULL,
  `timestp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `percurso` int NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Indexes for dumped tables
--

--
-- Indexes for table `athletes`
--
ALTER TABLE `athletes`
  ADD PRIMARY KEY (`num`),
  ADD KEY `event_id` (`event_id`),
  ADD KEY `track_id` (`track_id`);

--
-- Indexes for table `athletes_times`
--
ALTER TABLE `athletes_times`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `id` (`id`),
  ADD KEY `athlete_num` (`athlete_num`);

--
-- Indexes for table `checkpoints`
--
ALTER TABLE `checkpoints`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `classificacao`
--
ALTER TABLE `classificacao`
  ADD UNIQUE KEY `athlete_2` (`athlete`),
  ADD KEY `athlete` (`athlete`);

--
-- Indexes for table `equipamento`
--
ALTER TABLE `equipamento`
  ADD PRIMARY KEY (`id`),
  ADD KEY `event_id` (`event_id`),
  ADD KEY `checkpoint_id` (`checkpoint_id`);

--
-- Indexes for table `event_data`
--
ALTER TABLE `event_data`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `invalidos`
--
ALTER TABLE `invalidos`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `rede`
--
ALTER TABLE `rede`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `staffs`
--
ALTER TABLE `staffs`
  ADD PRIMARY KEY (`id`),
  ADD KEY `event_id` (`event_id`);

--
-- Indexes for table `tracks`
--
ALTER TABLE `tracks`
  ADD PRIMARY KEY (`id`),
  ADD KEY `event_id` (`event_id`);

--
-- Indexes for table `unsent_athletes`
--
ALTER TABLE `unsent_athletes`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `id` (`id`),
  ADD KEY `athlete_num` (`athlete_num`),
  ADD KEY `percursoid` (`percurso`),
  ADD KEY `percurso` (`percurso`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `athletes_times`
--
ALTER TABLE `athletes_times`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `invalidos`
--
ALTER TABLE `invalidos`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `rede`
--
ALTER TABLE `rede`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `unsent_athletes`
--
ALTER TABLE `unsent_athletes`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;

--
-- Constraints for dumped tables
--

--
-- Constraints for table `athletes`
--
ALTER TABLE `athletes`
  ADD CONSTRAINT `athletes_ibfk_1` FOREIGN KEY (`event_id`) REFERENCES `event_data` (`id`);

--
-- Constraints for table `athletes_times`
--
ALTER TABLE `athletes_times`
  ADD CONSTRAINT `athletes_times_ibfk_1` FOREIGN KEY (`athlete_num`) REFERENCES `athletes` (`num`);

--
-- Constraints for table `classificacao`
--
ALTER TABLE `classificacao`
  ADD CONSTRAINT `classificacao_ibfk_1` FOREIGN KEY (`athlete`) REFERENCES `athletes_times` (`athlete_num`);

--
-- Constraints for table `staffs`
--
ALTER TABLE `staffs`
  ADD CONSTRAINT `staffs_ibfk_1` FOREIGN KEY (`event_id`) REFERENCES `event_data` (`id`);

--
-- Constraints for table `tracks`
--
ALTER TABLE `tracks`
  ADD CONSTRAINT `tracks_ibfk_1` FOREIGN KEY (`event_id`) REFERENCES `event_data` (`id`);

--
-- Constraints for table `unsent_athletes`
--
ALTER TABLE `unsent_athletes`
  ADD CONSTRAINT `unsent_athletes_ibfk_1` FOREIGN KEY (`athlete_num`) REFERENCES `athletes` (`num`),
  ADD CONSTRAINT `unsent_athletes_ibfk_2` FOREIGN KEY (`percurso`) REFERENCES `tracks` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;

/*FIX_FILE_PERM*/
GRANT FILE on *.* to mytempo_user@'%';

