# Mailchecker

This Go script checks the email configuration of a domain by verifying the presence of MX, SPF, and DMARC records. It reads domain names from standard input and outputs the results for each domain. The script outputs a CSV header followed by the results, indicating whether the domain has MX records, SPF records (with the SPF configuration), and DMARC records (with the DMARC configuration).

## Usage

1. Compile the script:
   ```sh
   go build -o main
   ```

2. Run the compiled program and provide domain names as input:
   ```sh
   ./main
   ```
   Enter the domain names one per line and press `Ctrl+D` to end input.

## Output

The script prints the following CSV formatted fields:
```
domain, hasMX, hasSPF, spfRecords, hasDMARC, dmarcRecords
```

For each domain, it outputs:
- `domain`: The domain name.
- `hasMX`: Indicates if the domain has MX records (`true` or `false`).
- `hasSPF`: Indicates if the domain has SPF records (`true` or `false`).
- `spfRecords`: The SPF record(s) if present.
- `hasDMARC`: Indicates if the domain has DMARC records (`true` or `false`).
- `dmarcRecords`: The DMARC record(s) if present.

### Example

```sh
example.com, true, true, v=spf1 include:_spf.google.com ~all, true, v=DMARC1; p=none; rua=mailto:dmarc-reports@example.com
```
