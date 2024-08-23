-- phpMyAdmin SQL Dump
-- version 5.2.1
-- https://www.phpmyadmin.net/
--
-- Host: db
-- Generation Time: Aug 23, 2024 at 12:18 PM
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
  `name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `sex` varchar(1) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `team` varchar(80) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `track_id` int DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `athletes`
--

INSERT INTO `athletes` (`num`, `event_id`, `name`, `sex`, `team`, `track_id`) VALUES
(200, 201, 'josme merda recebA!', 'M', 'MERDA A JATO', 400);

-- --------------------------------------------------------

--
-- Table structure for table `athletes_times`
--

CREATE TABLE `athletes_times` (
  `antenna` int DEFAULT NULL,
  `athlete_num` int DEFAULT NULL,
  `athlete_time` varchar(12) COLLATE utf8mb4_unicode_ci DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `athletes_times`
--

INSERT INTO `athletes_times` (`antenna`, `athlete_num`, `athlete_time`) VALUES
(2, 200, '08:23:45.543');

-- --------------------------------------------------------

--
-- Table structure for table `event_data`
--

CREATE TABLE `event_data` (
  `id` int NOT NULL,
  `event_date` date DEFAULT NULL,
  `event_title` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `event_data`
--

INSERT INTO `event_data` (`id`, `event_date`, `event_title`) VALUES
(201, '2024-08-20', 'PROVA MUTO SEXO');

-- --------------------------------------------------------

--
-- Table structure for table `tracks`
--

CREATE TABLE `tracks` (
  `id` int NOT NULL,
  `event_id` int DEFAULT NULL,
  `race_description` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `km` decimal(3,1) DEFAULT NULL,
  `inicio` time DEFAULT NULL,
  `chegada` time DEFAULT NULL,
  `largada` time DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `tracks`
--

INSERT INTO `tracks` (`id`, `event_id`, `race_description`, `km`, `inicio`, `chegada`, `largada`) VALUES
(400, 201, 'BISTA', 5.0, '10:15:00', '10:20:00', '00:05:00');

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
  ADD KEY `athlete_num` (`athlete_num`);

--
-- Indexes for table `event_data`
--
ALTER TABLE `event_data`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `tracks`
--
ALTER TABLE `tracks`
  ADD PRIMARY KEY (`id`),
  ADD KEY `event_id` (`event_id`);

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
-- Constraints for table `tracks`
--
ALTER TABLE `tracks`
  ADD CONSTRAINT `tracks_ibfk_1` FOREIGN KEY (`event_id`) REFERENCES `event_data` (`id`);
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
