-- --------------------------------------------------------
-- Host:                         127.0.0.1
-- Server version:               PostgreSQL 15.14 (Debian 15.14-1.pgdg13+1) on x86_64-pc-linux-gnu, compiled by gcc (Debian 14.2.0-19) 14.2.0, 64-bit
-- Server OS:                    
-- HeidiSQL Version:             12.7.0.6850
-- --------------------------------------------------------

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES  */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

-- Dumping structure for table public.documents
CREATE TABLE IF NOT EXISTS "documents" (
	"id" SERIAL NOT NULL,
	"name" VARCHAR(500) NOT NULL,
	"path" TEXT NOT NULL,
	"content" TEXT NULL DEFAULT NULL,
	"uploaded" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	"page_size" INTEGER NULL DEFAULT 0,
	"title" VARCHAR(500) NULL DEFAULT NULL,
	"author" VARCHAR(500) NULL DEFAULT NULL,
	"subject" VARCHAR(500) NULL DEFAULT NULL,
	"keywords" TEXT NULL DEFAULT NULL,
	"creator" VARCHAR(500) NULL DEFAULT NULL,
	"producer" VARCHAR(500) NULL DEFAULT NULL,
	"is_encrypted" BOOLEAN NULL DEFAULT false,
	"pdf_version" VARCHAR(10) NULL DEFAULT NULL,
	"word_count" INTEGER NULL DEFAULT 0,
	"char_count" INTEGER NULL DEFAULT 0,
	"has_forms" BOOLEAN NULL DEFAULT false,
	"extraction_duration" BIGINT NULL DEFAULT 0,
	"extraction_success" BOOLEAN NULL DEFAULT false,
	"extraction_errors" JSONB NULL DEFAULT NULL,
	"file_size" BIGINT NULL DEFAULT 0,
	"file_hash" VARCHAR(64) NULL DEFAULT NULL,
	"mime_type" VARCHAR(100) NULL DEFAULT 'application/pdf',
	"created_at" TIMESTAMPTZ NULL DEFAULT CURRENT_TIMESTAMP,
	"updated_at" TIMESTAMPTZ NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY ("id")
);

-- Data exporting was unselected.

/*!40103 SET TIME_ZONE=IFNULL(@OLD_TIME_ZONE, 'system') */;
/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IFNULL(@OLD_FOREIGN_KEY_CHECKS, 1) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40111 SET SQL_NOTES=IFNULL(@OLD_SQL_NOTES, 1) */;
