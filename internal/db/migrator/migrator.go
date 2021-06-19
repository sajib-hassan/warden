package migrator

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/kamva/mgm/v3"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/sajib-hassan/warden/pkg/dbconn"
)

func ExecuteCreate(args []string) {
	startTime := time.Now()

	extPtr := viper.GetString("extPtr")
	dirPtr := viper.GetString("dirPtr")
	formatPtr := viper.GetString("formatPtr")
	timezoneName := viper.GetString("timezoneName")
	seq := viper.GetBool("seq")
	seqDigits := viper.GetInt("seqDigits")

	if len(args) == 0 {
		log.Fatal("error: please specify name")
	}

	name := args[0]

	if extPtr == "" {
		log.Fatal("error: --ext or -e flag must be specified")
	}

	timezone, err := time.LoadLocation(timezoneName)
	if err != nil {
		log.Fatal("error: ", err)
	}

	if err := createCmd(dirPtr, startTime.In(timezone), formatPtr, name, extPtr, seq, seqDigits, true); err != nil {
		log.Fatal("error: ", err)
	}
}

func ExecuteUp(args []string) {
	m, db := getMigrator()
	defer db.Client().Disconnect(mgm.Ctx())

	limit := -1
	if len(args) > 0 {
		n, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			log.Fatal("error: can't read limit argument N")
		}
		limit = int(n)
	}

	if err := upCmd(m, limit); err != nil {
		log.Fatal("error: ", err)
	}
}

func ExecuteDown(args []string) {
	m, db := getMigrator()
	defer db.Client().Disconnect(mgm.Ctx())

	num, needsConfirm, err := numDownMigrationsFromArgs(viper.GetBool("applyAll"), args)
	if err != nil {
		log.Fatal("error: ", err)
	}
	if needsConfirm {
		log.Println("Are you sure you want to apply all down migrations? [y/N]")
		var response string
		fmt.Scanln(&response)
		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" {
			log.Println("Applying all down migrations")
		} else {
			log.Fatal("Not applying all down migrations")
		}
	}

	if err := downCmd(m, num); err != nil {
		log.Fatal("error: ", err)
	}
}

func ExecuteDrop(args []string) {
	m, db := getMigrator()
	defer db.Client().Disconnect(mgm.Ctx())

	forceDrop := viper.GetBool("forceDrop")
	if !forceDrop {
		log.Println("Are you sure you want to drop the entire database schema? [y/N]")
		var response string
		fmt.Scanln(&response)
		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" {
			log.Println("Dropping the entire database schema")
		} else {
			log.Fatal("Aborted dropping the entire database schema")
		}
	}

	if err := dropCmd(m); err != nil {
		log.Fatal("error: ", err)
	}
}

func ExecuteForce(args []string) {
	m, db := getMigrator()
	defer db.Client().Disconnect(mgm.Ctx())

	if len(args) == 0 {
		log.Fatal("error: please specify version argument V")
	}

	v, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		log.Fatal("error: can't read version argument V")
	}

	if v < -1 {
		log.Fatal("error: argument V must be >= -1")
	}

	if err := forceCmd(m, int(v)); err != nil {
		log.Fatal("error: ", err)
	}
}

func ExecuteVersion(args []string) {
	m, db := getMigrator()
	defer db.Client().Disconnect(mgm.Ctx())

	if err := versionCmd(m); err != nil {
		log.Fatal("error: ", err)
	}
}

func ExecuteGoto(args []string) {
	m, db := getMigrator()
	defer db.Client().Disconnect(mgm.Ctx())

	if len(args) == 0 {
		log.Fatal("error: please specify version argument V")
	}

	v, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		log.Fatal("error: can't read version argument V")
	}

	if err := gotoCmd(m, uint(v)); err != nil {
		log.Fatal("error: ", err)
	}
}

func getMigrator() (*migrate.Migrate, *mongo.Database) {
	err := dbconn.Connect()
	if err != nil {
		log.Fatal("mongodb connect", err)
	}

	_, client, db, err := mgm.DefaultConfigs()
	if err != nil {
		log.Fatal("mongodb mgm get config", err)
	}

	driver, err := mongodb.WithInstance(client, &mongodb.Config{DatabaseName: db.Name()})
	if err != nil {
		log.Fatal("migrate mongodb instance: ", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		viper.GetString("source"),
		"mongodb",
		driver,
	)

	if err != nil {
		log.Fatal("migrate new with database instance ", err)
	}
	return m, db
}

func upCmd(m *migrate.Migrate, limit int) error {
	if limit >= 0 {
		if err := m.Steps(limit); err != nil {
			if err != migrate.ErrNoChange {
				return err
			}
			log.Println(err)
		}
	} else {
		if err := m.Up(); err != nil {
			if err != migrate.ErrNoChange {
				return err
			}
			log.Println(err)
		}
	}
	return nil
}

func downCmd(m *migrate.Migrate, limit int) error {
	if limit >= 0 {
		if err := m.Steps(-limit); err != nil {
			if err != migrate.ErrNoChange {
				return err
			}
			log.Println(err)
		}
	} else {
		if err := m.Down(); err != nil {
			if err != migrate.ErrNoChange {
				return err
			}
			log.Println(err)
		}
	}
	return nil
}

func dropCmd(m *migrate.Migrate) error {
	if err := m.Drop(); err != nil {
		return err
	}
	return nil
}

func forceCmd(m *migrate.Migrate, v int) error {
	if err := m.Force(v); err != nil {
		return err
	}
	return nil
}

func versionCmd(m *migrate.Migrate) error {
	v, dirty, err := m.Version()
	if err != nil {
		return err
	}
	if dirty {
		log.Printf("%v (dirty)\n", v)
	} else {
		log.Println(v)
	}
	return nil
}

func gotoCmd(m *migrate.Migrate, v uint) error {
	if err := m.Migrate(v); err != nil {
		if err != migrate.ErrNoChange {
			return err
		}
		log.Println(err)
	}
	return nil
}

// numDownMigrationsFromArgs returns an int for number of migrations to apply
// and a bool indicating if we need a confirm before applying
func numDownMigrationsFromArgs(applyAll bool, args []string) (int, bool, error) {
	if applyAll {
		if len(args) > 0 {
			return 0, false, errors.New("-all cannot be used with other arguments")
		}
		return -1, false, nil
	}

	switch len(args) {
	case 0:
		return -1, true, nil
	case 1:
		downValue := args[0]
		n, err := strconv.ParseUint(downValue, 10, 64)
		if err != nil {
			return 0, false, errors.New("can't read limit argument N")
		}
		return int(n), false, nil
	default:
		return 0, false, errors.New("too many arguments")
	}
}

const (
	defaultTimeFormat = "20060102150405"
)

var (
	errInvalidSequenceWidth     = errors.New("Digits must be positive")
	errIncompatibleSeqAndFormat = errors.New("The seq and format options are mutually exclusive")
	errInvalidTimeFormat        = errors.New("Time format may not be empty")
)

func nextSeqVersion(matches []string, seqDigits int) (string, error) {
	if seqDigits <= 0 {
		return "", errInvalidSequenceWidth
	}

	nextSeq := uint64(1)

	if len(matches) > 0 {
		filename := matches[len(matches)-1]
		matchSeqStr := filepath.Base(filename)
		idx := strings.Index(matchSeqStr, "_")

		if idx < 1 { // Using 1 instead of 0 since there should be at least 1 digit
			return "", fmt.Errorf("Malformed migration filename: %s", filename)
		}

		var err error
		matchSeqStr = matchSeqStr[0:idx]
		nextSeq, err = strconv.ParseUint(matchSeqStr, 10, 64)

		if err != nil {
			return "", err
		}

		nextSeq++
	}

	version := fmt.Sprintf("%0[2]*[1]d", nextSeq, seqDigits)

	if len(version) > seqDigits {
		return "", fmt.Errorf("Next sequence number %s too large. At most %d digits are allowed", version, seqDigits)
	}

	return version, nil
}

func timeVersion(startTime time.Time, format string) (version string, err error) {
	switch format {
	case "":
		err = errInvalidTimeFormat
	case "unix":
		version = strconv.FormatInt(startTime.Unix(), 10)
	case "unixNano":
		version = strconv.FormatInt(startTime.UnixNano(), 10)
	default:
		version = startTime.Format(format)
	}

	return
}

// createCmd (meant to be called via a CLI command) creates a new migration
func createCmd(dir string, startTime time.Time, format string, name string, ext string, seq bool, seqDigits int, print bool) error {
	if seq && format != defaultTimeFormat {
		return errIncompatibleSeqAndFormat
	}

	var version string
	var err error

	dir = filepath.Clean(dir)
	ext = "." + strings.TrimPrefix(ext, ".")

	if seq {
		matches, err := filepath.Glob(filepath.Join(dir, "*"+ext))

		if err != nil {
			return err
		}

		version, err = nextSeqVersion(matches, seqDigits)

		if err != nil {
			return err
		}
	} else {
		version, err = timeVersion(startTime, format)

		if err != nil {
			return err
		}
	}

	versionGlob := filepath.Join(dir, version+"_*"+ext)
	matches, err := filepath.Glob(versionGlob)

	if err != nil {
		return err
	}

	if len(matches) > 0 {
		return fmt.Errorf("duplicate migration version: %s", version)
	}

	if err = os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	for _, direction := range []string{"up", "down"} {
		basename := fmt.Sprintf("%s_%s.%s%s", version, name, direction, ext)
		filename := filepath.Join(dir, basename)

		if err = createFile(filename); err != nil {
			return err
		}

		if print {
			absPath, _ := filepath.Abs(filename)
			log.Println(absPath)
		}
	}

	return nil
}

func createFile(filename string) error {
	// create exclusive (fails if file already exists)
	// os.Create() specifies 0666 as the FileMode, so we're doing the same
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)

	if err != nil {
		return err
	}

	return f.Close()
}
