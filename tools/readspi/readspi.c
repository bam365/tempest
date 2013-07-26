/* spitest.c - A test to read an SPI signal from GPIO port 
 * Blake Mitchell, 2013
 */

#include <stdio.h>
#include <unistd.h>
#include <bcm2835.h>



#define PIN_CLK 	RPI_GPIO_P1_12
#define PIN_MISO	RPI_GPIO_P1_16
#define PIN_MOSI	RPI_GPIO_P1_18
#define PIN_CS		RPI_GPIO_P1_22


void setup_pins();
void prep_for_read(uint8_t adc);
int read_channel(uint8_t adc);


int main(int argc, char **argv)
{
	int chan;

	if (!bcm2835_init())
		return -1;	

	if ((chan = get_channel(argc, argv)) < 0) {
		fprintf(stderr, "Invalid channel argument\n");
		return -1;
	}

	setup_pins();
	bcm2835_delay(50);
	printf("%d\n", read_channel((uint8_t)chan));


	return 0;
}


int get_channel(int argc, char **argv)
{
	int ret = -1, buf;

	if (argc > 1) 
	if (sscanf(argv[1], "%d", &buf) == 1) 
	if (buf >= 0 && buf <= 7)
		ret = buf;
	
	return ret;
}	
		
		
void setup_pins()
{
	bcm2835_gpio_fsel(PIN_CLK, BCM2835_GPIO_FSEL_OUTP);
	bcm2835_gpio_fsel(PIN_MISO, BCM2835_GPIO_FSEL_INPT);
	bcm2835_gpio_set_pud(PIN_MISO, BCM2835_GPIO_PUD_UP);
	bcm2835_gpio_fsel(PIN_MOSI, BCM2835_GPIO_FSEL_OUTP);
	bcm2835_gpio_fsel(PIN_CS, BCM2835_GPIO_FSEL_OUTP);
}


void cycle_clock()
{
	bcm2835_gpio_write(PIN_CLK, HIGH);
	bcm2835_delay(5);
	bcm2835_gpio_write(PIN_CLK, LOW);
}


void prep_for_read(uint8_t adc)
{
	uint8_t cmdbyte, out;
	int i;

	cmdbyte = (adc | 0x18) << 3;
	for (i = 0; i < 5; i++) {
		bcm2835_gpio_write(PIN_MOSI, (cmdbyte & 0x80) ? HIGH : LOW);
		cmdbyte <<= 1;
		cycle_clock();
	}
}
		
		

int read_channel(uint8_t adc)
{
	int inp, i;

	inp = 0;
	bcm2835_gpio_write(PIN_CS, HIGH);
	bcm2835_gpio_write(PIN_CLK, LOW);
	bcm2835_gpio_write(PIN_CS, LOW);
	prep_for_read(adc);
	/* One empty bit, one NUL bit, 10 data bytes */
	for (i = 0; i < 12; i++) {
		cycle_clock();
		inp <<= 1;
		inp |= bcm2835_gpio_lev(PIN_MISO);
	}
	bcm2835_gpio_write(PIN_CS, HIGH);

	/* Drop NUL bit */
	inp >>= 1;
	/* Mask out all but 10 bits */
	inp &= 0x000003FF;
	return inp;
}
		
		
		


	

	
