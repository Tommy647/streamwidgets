
streamwidget: clean
	@echo 'building'; go build -o streamwidget ./cmd/streamwidget

clean:
	@echo 'cleaning'; rm -Rf streamwidget
