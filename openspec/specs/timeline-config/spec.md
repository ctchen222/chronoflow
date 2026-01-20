# timeline-config Specification

## Purpose
TBD - created by archiving change enhance-day-view-timeline. Update Purpose after archive.
## Requirements
### Requirement: Timeline Time Range Configuration

Users MUST be able to configure the visible time range of the timeline.

#### Scenario: Default time range

- **Given** no custom configuration exists
- **When** the Day View timeline renders
- **Then** the timeline shows times from 08:00 to 18:00

#### Scenario: Custom time range

- **Given** config has day_start "06:00" and day_end "22:00"
- **When** the Day View timeline renders
- **Then** the timeline shows times from 06:00 to 22:00

### Requirement: Timeline Slot Granularity Configuration

Users SHALL be able to configure the minimum time slot size.

#### Scenario: Default slot granularity

- **Given** no custom configuration exists
- **When** the timeline renders
- **Then** time markers appear at 30-minute intervals

#### Scenario: Custom slot granularity

- **Given** config has slot_minutes 15
- **When** the timeline renders
- **Then** time markers appear at 15-minute intervals

### Requirement: Configuration Persistence

Timeline settings MUST be persisted to a configuration file.

#### Scenario: Config file location

- **Given** user modifies timeline settings
- **When** settings are saved
- **Then** settings are stored in ~/.chronoflow/config.json

#### Scenario: Config file auto-creation

- **Given** ~/.chronoflow/config.json does not exist
- **When** the application starts
- **Then** a config file is created with default values

#### Scenario: Config loaded on startup

- **Given** ~/.chronoflow/config.json exists with custom values
- **When** the application starts
- **Then** timeline uses the configured values

