library(dplyr)

pause <- function() {
  if (interactive())
  {
    invisible(readline(prompt = "Press <Enter> to continue..."))
  }
  else
  {
    cat("Press <Enter> to continue...")
    invisible(readLines(file("stdin"), 1))
  }
}

args = commandArgs(trailingOnly=TRUE)
if (length(args) < 1) {
  print("Must provide a results file!")
  quit()
}
res <- read.csv(args[1]);
res$s <- 0;

for (i in 1:nrow(res)) {
  cur <- res[i,];
  if (cur$n == 1) { 
    base <- cur$t; 
  }
  cur$s <- base / cur$t;
  res[i,] <- cur;
}

print(res)
pause()


drawGraph <- function(column, name_filter, title, ylab, show_one = FALSE, hide_omp = FALSE) {
  to_add <- FALSE
  colors <- c("red", "green", "blue")
  leg <- c()
  j <- 1;
  dataset <- rev(unique(res$aname[grepl(name_filter, res$aname)]));
  if (hide_omp) {
    dataset <- dataset[!grepl("omp", dataset)];
  }
  print(sort(filter(res, aname %in% dataset)[[column]]))
  for (i in dataset) {
    if (show_one) {
      filt <- filter(filter(res, aname == i), s == 1 | m == 100);
      yran <- range(filter(res, (s == 1 | m == 100) & aname %in% dataset)[[column]])
      x <- 1:(length(filt[[column]]))
      print(x)
    } else {
      filt <- filter(filter(res, aname == i), s > 1 | m == 100);
      yran <- range(filter(res, (s > 1 | m == 100) & aname %in% dataset)[[column]])
      x <- filt$n
    }

    leg <- unlist(c(leg, i))
    y <- filt[[column]]
    if (to_add) {
      lines(x, y, col=colors[j], type="b")
    } else {
      plot(x=x, y=y, ylim=yran, xlab = "N. of processors", ylab = ylab, main = title, type="b", col=colors[j])
    }
    to_add <- TRUE
    j <- j + 1
  }
  print(par("col"));
  legend("topright", col=colors, legend=leg, pch=15);
}

# Draw time chart for strong scalability
png("Strong_Time.png", 
  width     = 3.25,
  height    = 3.25,
  units     = "in",
  res       = 1200,
  pointsize = 4)
drawGraph("t", "strong", "Strong scalability", "Execution time");
# Draw speedup chart for strong scalability
png("Strong_Speedup.png", 
  width     = 3.25,
  height    = 3.25,
  units     = "in",
  res       = 1200,
  pointsize = 4)

drawGraph("s", "strong", "Strong scalability", "Speedup");
# Draw time chart for weak scalability
png("Weak_Time.png", 
  width     = 3.25,
  height    = 3.25,
  units     = "in",
  res       = 1200,
  pointsize = 4)

drawGraph("t", "weak", "Weak scalability", "Execution time");
# Draw speedup chart for weak scalability
png("Weak_Speedup.png", 
  width     = 3.25,
  height    = 3.25,
  units     = "in",
  res       = 1200,
  pointsize = 4)
drawGraph("s", "weak", "Weak scalability", "Scaled speedup");

# Draw time chart for strong scalability
png("Strong_Time_MPI_Only.png", 
  width     = 3.25,
  height    = 3.25,
  units     = "in",
  res       = 1200,
  pointsize = 4)
drawGraph("t", "strong", "Strong scalability", "Execution time", FALSE, TRUE);
# Draw speedup chart for strong scalability
png("Strong_Speedup_MPI_Only.png", 
  width     = 3.25,
  height    = 3.25,
  units     = "in",
  res       = 1200,
  pointsize = 4)

drawGraph("s", "strong", "Strong scalability", "Speedup", FALSE, TRUE);
# Draw time chart for weak scalability
png("Weak_Time_MPI_Only.png", 
  width     = 3.25,
  height    = 3.25,
  units     = "in",
  res       = 1200,
  pointsize = 4)

drawGraph("t", "weak", "Weak scalability", "Execution time", FALSE, TRUE);
# Draw speedup chart for weak scalability
png("Weak_Speedup_MPI_Only.png", 
  width     = 3.25,
  height    = 3.25,
  units     = "in",
  res       = 1200,
  pointsize = 4)
drawGraph("s", "weak", "Weak scalability", "Scaled speedup", FALSE, TRUE);

#warnings()
